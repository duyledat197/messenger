package search

import (
	"context"

	index "github.com/blevesearch/bleve_index_api"
)

// Search represents the capabilities of a search engine.
type Search interface {
	// Index adds a document to the search engine.
	//
	// It takes a context.Context for cancellation and a string ID and data of any type.
	// It returns an error if the operation fails.
	Index(ctx context.Context, id string, data any) error

	// Delete removes a document from the search engine.
	//
	// It takes a context.Context for cancellation and a string ID.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id string) error

	// Document retrieves a document from the search engine.
	//
	// It takes a context.Context for cancellation and a string ID.
	// It returns a bleve_index_api.Document and an error.
	Document(ctx context.Context, id string) (index.Document, error)

	// DocCount returns the total number of documents in the search engine.
	//
	// It takes a context.Context for cancellation.
	// It returns a uint64 and an error.
	DocCount(ctx context.Context) (uint64, error)

	// Search performs a query on the search engine and returns the results.
	//
	// It takes a context.Context for cancellation and a filter of any type.
	// It returns any type (the result of the search) and an error.
	Search(ctx context.Context, filter any) (any, error)

	// Fields retrieves the list of fields available in the search engine.
	//
	// It returns a slice of strings and an error.
	Fields() ([]string, error)
}
