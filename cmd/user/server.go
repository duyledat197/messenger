package cmd

import (
	"context"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/user/repository"
	"openmyth/messgener/internal/user/repository/postgres"
	"openmyth/messgener/internal/user/service"
	pb "openmyth/messgener/pb/user"
	"openmyth/messgener/pkg/grpc_server"
	"openmyth/messgener/pkg/postgres_client"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/util/snowflake"
)

var server struct {
	// database
	pgClient *postgres_client.PostgresClient

	// config
	cfg *config.Config

	idGenerator snowflake.Generator

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

// loadConfigs loads the configuration for the server.
func loadConfigs() {
	server.cfg = config.LoadConfig()

	server.idGenerator = *snowflake.NewGenerator(1)
}

// loadDatabases initializes the database clients for the server.
func loadDatabases() {
	cfg := server.cfg
	server.pgClient = postgres_client.NewPostgresClient(cfg.User.Postgres.Address())

	server.lifecycle.WithFactories(
		server.pgClient,
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
	cfg := server.cfg

	srv := grpc_server.NewGrpcServer(cfg.User.Endpoint)

	pb.RegisterAuthServiceServer(srv, server.authService)

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
