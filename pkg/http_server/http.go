package http_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	"openmyth/messgener/config"
)

type HttpServer struct {
	mux      *runtime.ServeMux
	server   *http.Server
	endpoint *config.Endpoint
}

// NewHttpServer creates a new HTTP server with the provided handler and endpoint.
func NewHttpServer(
	handler func(mux *runtime.ServeMux),
	endpoint *config.Endpoint,
) *HttpServer {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{UseEnumNumbers: false, EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{AllowPartial: true},
		}),
		runtime.WithMetadata(MapMetaDataWithBearerToken()),
		runtime.WithErrorHandler(forwardErrorResponse),
	)
	handler(mux)
	middlewares := []middlewareFunc{
		allowCORS,
	}

	slices.Reverse(middlewares)

	var handleR http.Handler = mux
	for _, handle := range middlewares {
		handleR = handle(handleR)
	}

	return &HttpServer{
		mux:      mux,
		endpoint: endpoint,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", endpoint.Port),
			Handler: handleR,
		},
	}
}

// Start starts the HTTP server.
func (s *HttpServer) Start(_ context.Context) error {
	log.Printf("Server listin in port: %s\n", s.endpoint.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Stop stops the HTTP server gracefully.
func (s *HttpServer) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
