package cmd

import (
	"context"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/internal/chat/repository/cqlx"
	open_search "openmyth/messgener/internal/chat/repository/opensearch"
	redisRepo "openmyth/messgener/internal/chat/repository/redis"
	"openmyth/messgener/internal/chat/service"
	pb "openmyth/messgener/pb/chat"
	"openmyth/messgener/pkg/courier"
	"openmyth/messgener/pkg/grpc_server"
	"openmyth/messgener/pkg/logger"
	"openmyth/messgener/pkg/opensearch"
	"openmyth/messgener/pkg/postgres_client"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/pkg/redis"
	"openmyth/messgener/pkg/scylla"
	"openmyth/messgener/util/snowflake"
)

var server struct {
	// database
	pgClient         *postgres_client.PostgresClient
	openSearchClient *opensearch.Client
	scylladbClient   *scylla.ScyllaClient
	courierClient    *courier.Client
	redisClient      *redis.Client

	// repository
	messageRepo      repository.MessageRepository
	channelRepo      repository.ChannelRepository
	cacheChannelRepo repository.CacheChannelRepository

	// service
	messageService pb.MessageServiceServer
	channelService pb.ChannelServiceServer

	lifecycle *processor.Lifecycle
}

func loadLogger() {
	logger.SetLoggerGlobal()
}

// loadLifecycle initializes the server's lifecycle by creating a new instance of the Lifecycle struct.
func loadLifecycle() {
	server.lifecycle = processor.NewLifecycle()
}

// loadConfigs loads the configuration for the server.
func loadConfigs() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
}

func loadIDGenerator() {
	num, _ := server.redisClient.Client.Incr(context.Background(), "node_num").Result()
	snowflake.SetGlobalIDGenerator(num)
}

// loadDatabases initializes the database clients for the server.
func loadDatabases() {
	server.openSearchClient = opensearch.NewOpenSearch(config.GetGlobalConfig().Chat.OpenSearch)
	server.scylladbClient = scylla.NewScylla(config.GetGlobalConfig().Chat.ScyllaDB)
	// server.courierClient = courier.NewClient(config.GetGlobalConfig().Chat.Courier)
	server.redisClient = redis.NewClient(config.GetGlobalConfig().Chat.Redis)
	server.pgClient = postgres_client.NewPostgresClient(config.GetGlobalConfig().Chat.Postgres)

	server.lifecycle.WithFactories(
		server.pgClient,
		server.openSearchClient,
		server.redisClient,
		server.scylladbClient,
		// server.courierClient,
	)
}

// loadRepositories initializes the message and channel repositories for the server.
func loadRepositories() {
	server.messageRepo = cqlx.NewMessageRepository(server.scylladbClient)
	server.channelRepo = open_search.NewChannelRepository(server.openSearchClient)
	server.cacheChannelRepo = redisRepo.NewCacheChannelRepository(server.redisClient)
}

// loadServices initializes the message and channel services for the server.
func loadServices() {
	server.messageService = service.NewMessageService(
		server.messageRepo,
		server.channelRepo)

	server.channelService = service.NewChannelService(server.channelRepo, server.cacheChannelRepo)
}

func loadServer() {

	srv := grpc_server.NewGrpcServer(config.GetGlobalConfig().Chat.Endpoint)

	pb.RegisterChannelServiceServer(srv, server.channelService)
	pb.RegisterMessageServiceServer(srv, server.messageService)

	server.lifecycle.WithProcessors(srv)
}

// Load initializes the server with the necessary configurations and dependencies.
func Load() {
	loadLifecycle()
	loadConfigs()
	loadLogger()
	loadDatabases()
	loadIDGenerator()
	loadRepositories()
	loadServices()
	loadServer()
}

// Start starts the server with the given context.
func Start(ctx context.Context) {
	server.lifecycle.Start(ctx)
}
