package cmd

import (
	"context"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/internal/chat/repository/cqlx"
	open_search "openmyth/messgener/internal/chat/repository/opensearch"
	"openmyth/messgener/internal/chat/service"
	pb "openmyth/messgener/pb/chat"
	"openmyth/messgener/pkg/courier"
	"openmyth/messgener/pkg/grpc_server"
	"openmyth/messgener/pkg/opensearch"
	"openmyth/messgener/pkg/postgres_client"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/pkg/scylla"
	"openmyth/messgener/util/snowflake"
)

var server struct {
	// database
	pgClient         *postgres_client.PostgresClient
	openSearchClient *opensearch.Client
	scylladbClient   *scylla.ScyllaClient
	courierClient    *courier.Client

	// config
	cfg *config.Config

	idGenerator snowflake.Generator

	// repository
	messageRepo repository.MessageRepository
	channelRepo repository.ChannelRepository

	// service
	messageService pb.MessageServiceServer
	channelService pb.ChannelServiceServer

	lifecycle *processor.Lifecycle
}

// loadLifecycle initializes the server's lifecycle by creating a new instance of the Lifecycle struct.
func loadLifecycle() {
	server.lifecycle = processor.NewLifecycle()
}

// loadConfigs loads the configuration for the server.
func loadConfigs() {
	server.cfg = config.LoadConfig()

	server.idGenerator = *snowflake.NewGenerator(1)
}

// loadDatabases initializes the database clients for the server.
func loadDatabases() {
	cfg := server.cfg
	// server.pgClient = postgres_client.NewPostgresClient(cfg.PostgresDB.Address())
	server.openSearchClient = opensearch.NewOpenSearch(cfg.OpenSearch)
	server.scylladbClient = scylla.NewScylla(cfg.ScyllaDB)
	server.courierClient = courier.NewClient(cfg.Courier)

	server.lifecycle.WithFactories(
		// server.pgClient,
		server.openSearchClient, server.scylladbClient, server.courierClient)
}

// loadRepositories initializes the message and channel repositories for the server.
func loadRepositories() {
	server.messageRepo = cqlx.NewMessageRepository(server.scylladbClient)
	server.channelRepo = open_search.NewChannelRepository(server.openSearchClient)
}

// loadServices initializes the message and channel services for the server.
func loadServices() {
	server.messageService = service.NewMessageService(
		server.idGenerator, server.messageRepo,
		server.channelRepo)

	server.channelService = service.NewChannelService(server.channelRepo, server.idGenerator)

}

func loadServer() {
	cfg := server.cfg

	srv := grpc_server.NewGrpcServer(cfg.ChatService)

	pb.RegisterChannelServiceServer(srv, server.channelService)
	pb.RegisterMessageServiceServer(srv, server.messageService)

	server.lifecycle.WithProcessors(srv)
}

// Load initializes the server with the necessary configurations and dependencies.
func Load() {
	loadLifecycle()
	loadConfigs()
	loadDatabases()
	loadRepositories()
	loadServices()
	loadServer()
}

// Start starts the server with the given context.
func Start(ctx context.Context) {
	server.lifecycle.Start(ctx)
}
