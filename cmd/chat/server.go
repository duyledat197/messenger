package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"openmyth/messgener/config"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/internal/chat/repository/cqlx"
	open_search "openmyth/messgener/internal/chat/repository/opensearch"
	"openmyth/messgener/internal/chat/service"
	pb "openmyth/messgener/pb/chat"
	"openmyth/messgener/pkg/courier"
	"openmyth/messgener/pkg/grpc_server"
	"openmyth/messgener/pkg/opensearch"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/pkg/scylla"
	"openmyth/messgener/util/snowflake"
)

var server struct {
	// database
	// pgClient         *postgres_client.PostgresClient
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

	factories  []processor.Factory
	processors []processor.Processor
}

// Load initializes the server with the necessary configurations and dependencies.
func Load() {
	cfg := config.LoadConfig()

	server.idGenerator = *snowflake.NewGenerator(1)

	server.openSearchClient = opensearch.NewOpenSearch(cfg.OpenSearch)
	server.scylladbClient = scylla.NewScylla(cfg.ScyllaDB)
	server.courierClient = courier.NewClient(cfg.Courier)

	server.factories = append(server.factories, server.openSearchClient, server.scylladbClient, server.courierClient)

	server.messageRepo = cqlx.NewMessageRepository(server.scylladbClient)
	server.channelRepo = open_search.NewChannelRepository(server.openSearchClient)

	server.messageService = service.NewMessageService(
		server.idGenerator, server.messageRepo,
		server.channelRepo)

	server.channelService = service.NewChannelService(server.channelRepo, server.idGenerator)

	// clients

	srv := grpc_server.NewGrpcServer(cfg.ChatService)

	pb.RegisterChannelServiceServer(srv, server.channelService)
	pb.RegisterMessageServiceServer(srv, server.messageService)

	server.processors = append(server.processors, srv)
}

// Start starts the server and initializes the necessary components.
//
// It connects to the factories and processors, and starts them in separate goroutines.
// It listens for interrupt or termination signals and calls the stop function accordingly.
func Start(ctx context.Context) {
	errChan := make(chan error)
	signChan := make(chan os.Signal, 1)

	for _, f := range server.factories {
		if err := f.Connect(ctx); err != nil {
			errChan <- err
		}
	}

	for _, p := range server.processors {
		go func(p processor.Processor) {
			if err := p.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}

	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	select {
	case _ = <-errChan:
		stop(ctx)
	case <-signChan:
		stop(ctx, true)
	}

}

// stop stops the server gracefully by closing all factories and starting all processors.
func stop(ctx context.Context, graceful ...bool) {
	for _, p := range server.processors {
		if err := p.Stop(ctx); err != nil {
			slog.Error("unable to close processor:", err)
		}
	}

	if len(graceful) > 0 {
		time.Sleep(5 * time.Second)
	}

	for _, f := range server.factories {
		if err := f.Close(ctx); err != nil {
			slog.Error("unable to close factory:", err)
		}
	}

}
