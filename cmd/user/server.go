package cmd

import (
	"context"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/user/repository"
	"openmyth/messgener/internal/user/repository/postgres"
	"openmyth/messgener/internal/user/service"
	pb "openmyth/messgener/pb/user"
	"openmyth/messgener/pkg/grpc_server"
	"openmyth/messgener/pkg/logger"
	"openmyth/messgener/pkg/postgres_client"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/pkg/redis"
	"openmyth/messgener/util/snowflake"
)

var server struct {
	// database
	pgClient    *postgres_client.PostgresClient
	redisClient *redis.Client

	// repository
	userRepo repository.UserRepository

	// service
	authService pb.AuthServiceServer

	lifecycle *processor.Lifecycle
}

// loadLifecycle initializes the server's lifecycle by creating a new instance of the Lifecycle struct.
func loadLifecycle() {
	server.lifecycle = processor.NewLifecycle()
}

func loadLogger() {
	logger.SetLoggerGlobal()
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
	server.pgClient = postgres_client.NewPostgresClient(config.GetGlobalConfig().User.Postgres)
	server.redisClient = redis.NewClient(config.GetGlobalConfig().Chat.Redis)

	server.lifecycle.WithFactories(
		server.pgClient,
		server.redisClient,
	)
}

// loadRepositories initializes the message and channel repositories for the server.
func loadRepositories() {
	server.userRepo = postgres.NewUserRepository()
}

// loadServices initializes the message and channel services for the server.
func loadServices() {
	server.authService = service.NewAuthService(server.pgClient, server.userRepo)
}

func loadServer() {
	srv := grpc_server.NewGrpcServer(config.GetGlobalConfig().User.Endpoint)

	pb.RegisterAuthServiceServer(srv, server.authService)

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
