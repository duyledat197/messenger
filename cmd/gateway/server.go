package cmd

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"openmyth/messgener/config"
	chatPb "openmyth/messgener/pb/chat"
	userPb "openmyth/messgener/pb/user"
	"openmyth/messgener/pkg/grpc_client"
	"openmyth/messgener/pkg/http_server"
	"openmyth/messgener/pkg/processor"
	"openmyth/messgener/pkg/websocket"
)

var server struct {
	lifecycle *processor.Lifecycle

	messageClient chatPb.MessageServiceClient
	channelClient chatPb.ChannelServiceClient
	authClient    userPb.AuthServiceClient

	engine *websocket.Engine[chatPb.Message]
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

// loadServices initializes the message and channel services for the server.
func loadClients() {
	userConn := grpc_client.NewGrpcClient(config.GetGlobalConfig().User.Endpoint)
	chatConn := grpc_client.NewGrpcClient(config.GetGlobalConfig().Chat.Endpoint)

	server.authClient = userPb.NewAuthServiceClient(userConn)
	server.messageClient = chatPb.NewMessageServiceClient(chatConn)
	server.channelClient = chatPb.NewChannelServiceClient(chatConn)

	server.lifecycle.WithFactories(userConn, chatConn)
}

func loadEngine() {
	server.engine = websocket.NewEngine[chatPb.Message]()

	server.lifecycle.WithProcessors(server.engine)
}

func loadServer() {
	ctx := context.Background()
	srv := http_server.NewHttpServer(func(mux *runtime.ServeMux) {
		userPb.RegisterAuthServiceHandlerClient(ctx, mux, server.authClient)

		chatPb.RegisterChannelServiceHandlerClient(ctx, mux, server.channelClient)
		chatPb.RegisterMessageServiceHandlerClient(ctx, mux, server.messageClient)

		mux.HandlePath(http.MethodPost, "/ws", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			websocket.ServeWs(server.engine, w, r)
		})

	}, config.GetGlobalConfig().Gateway.Endpoint)

	server.lifecycle.WithProcessors(srv)
}

// Load initializes the server with the necessary configurations and dependencies.
func Load() {
	loadLifecycle()
	loadConfigs()
	loadClients()
	loadEngine()
	loadServer()
}

// Start starts the server with the given context.
func Start(ctx context.Context) {
	server.lifecycle.Start(ctx)
}
