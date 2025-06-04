package cmd

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"openmyth/messgener/config"
	chatPb "openmyth/messgener/pb/chat"
	userPb "openmyth/messgener/pb/user"
	"openmyth/messgener/pkg/grpc_client"
	"openmyth/messgener/pkg/http_server"
	"openmyth/messgener/pkg/processor"
)

var server struct {

	// config
	cfg *config.Config

	lifecycle *processor.Lifecycle

	messageClient chatPb.MessageServiceClient
	channelClient chatPb.ChannelServiceClient
	authClient    userPb.AuthServiceClient
}

// loadLifecycle initializes the server's lifecycle by creating a new instance of the Lifecycle struct.
func loadLifecycle() {
	server.lifecycle = processor.NewLifecycle()
}

// loadConfigs loads the configuration for the server.
func loadConfigs() {
	server.cfg = config.LoadConfig()
}

// loadServices initializes the message and channel services for the server.
func loadClients() {
	userConn := grpc_client.NewGrpcClient(server.cfg.User.Endpoint)
	chatConn := grpc_client.NewGrpcClient(server.cfg.Chat.Endpoint)

	server.authClient = userPb.NewAuthServiceClient(userConn)
	server.messageClient = chatPb.NewMessageServiceClient(chatConn)
	server.channelClient = chatPb.NewChannelServiceClient(chatConn)

	server.lifecycle.WithFactories(userConn, chatConn)
}

func loadServer() {
	ctx := context.Background()
	srv := http_server.NewHttpServer(func(mux *runtime.ServeMux) {
		userPb.RegisterAuthServiceHandlerClient(ctx, mux, server.authClient)

		chatPb.RegisterChannelServiceHandlerClient(ctx, mux, server.channelClient)
		chatPb.RegisterMessageServiceHandlerClient(ctx, mux, server.messageClient)
	}, server.cfg.Gateway.Endpoint)

	server.lifecycle.WithProcessors(srv)
}

// Load initializes the server with the necessary configurations and dependencies.
func Load() {
	loadLifecycle()
	loadConfigs()
	loadClients()
	loadServer()
}

// Start starts the server with the given context.
func Start(ctx context.Context) {
	server.lifecycle.Start(ctx)
}
