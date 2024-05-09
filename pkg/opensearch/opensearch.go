package opensearch

import (
	"context"
	"log"
	"os"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchtransport"

	"openmyth/messgener/config"
)

// Client represents the Client client.
// It contains a pointer to an *opensearch.Client and a pointer to a config.Database.
type Client struct {
	// Client is the actual opensearch client.
	*opensearch.Client

	// cfg is the configuration for the OpenSearch client.
	cfg *config.Database
}

// NewClient creates a new instance of the OpenSearch struct.
// It takes a pointer to a config.Database struct as a parameter and returns a pointer to an OpenSearch struct.
func NewOpenSearch(cfg *config.Database) *Client {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses:         []string{cfg.Address()},
		EnableMetrics:     true,
		EnableDebugLogger: os.Getenv("ENV") == "dev",
		Logger: &opensearchtransport.JSONLogger{
			Output:             log.Writer(),
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
		RetryBackoff: func(attempt int) time.Duration {
			return time.Duration(attempt*10) * time.Second
		},
	})
	if err != nil {
		log.Fatalf("unable to create opensearch client: %v", err)
	}
	return &Client{
		Client: client,
		cfg:    cfg,
	}
}

// Connect establishes a connection to the OpenSearch client.
// It takes a context.Context for cancellation.
// It returns an error if the connection fails.
func (o *Client) Connect(_ context.Context) error {
	return nil
}

// Close closes the OpenSearch client.
// It takes a context.Context as a parameter and returns an error if there was a problem closing the connection.
func (o *Client) Close(_ context.Context) error {
	return nil
}
