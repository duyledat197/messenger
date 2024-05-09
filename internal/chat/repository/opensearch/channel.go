package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/opensearch"
)

type channelRepository struct {
	opensearchClient *opensearch.Client
}

// NewChannelRepository creates a new instance of the ChannelRepository interface.
// It takes an opensearchClient of type *opensearch.Client as a parameter and returns a pointer to a channelRepository struct that implements the ChannelRepository interface.
func NewChannelRepository(opensearchClient *opensearch.Client) repository.ChannelRepository {
	return &channelRepository{opensearchClient}
}

// Create creates a new channel in the channel repository.
func (r *channelRepository) Create(ctx context.Context, e *entity.Channel) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	document := bytes.NewReader(b)

	req := opensearchapi.IndexRequest{
		Index:      e.TableName(),
		DocumentID: fmt.Sprint(e.ChannelID),
		Body:       document,
	}
	resp, err := req.Do(ctx, r.opensearchClient)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("create error: %s", resp.String())
	}

	return nil
}

// RetrieveByChannelID retrieves a channel by its ID from the repository.
func (r *channelRepository) RetrieveByChannelID(ctx context.Context, id int64) (*entity.Channel, error) {
	content := strings.NewReader(fmt.Sprintf(`{
		"from": 0,
		"size": 1,
		"query": {
			"term": {
				"_id": %d
			}
		}
	}
}`, id))
	e := &entity.Channel{}
	search := opensearchapi.SearchRequest{
		Index: []string{e.TableName()},
		Body:  content,
	}

	searchResponse, err := search.Do(ctx, r.opensearchClient)
	if err != nil {
		return nil, err
	}
	if searchResponse.IsError() {
		return nil, fmt.Errorf("retrieve error: %s", searchResponse.String())
	}

	defer searchResponse.Body.Close()

	dataB, err := io.ReadAll(searchResponse.Body)
	if err != nil {
		return nil, err
	}

	result := new(entity.Channel)
	if err := json.Unmarshal(dataB, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SearchByName searches for channels by name in the repository.
// It takes a context.Context and the name to search for as parameters.
// It returns a slice of *entity.Channel and an error.
func (r *channelRepository) SearchByName(ctx context.Context, name string) ([]*entity.Channel, error) {
	content := strings.NewReader(fmt.Sprintf(`{
		"from": 0,
		"size": 10,
		"query": {
			"match": {
				"name": %s
			}
		}
	}
}`, name))
	e := &entity.Channel{}
	search := opensearchapi.SearchRequest{
		Index: []string{e.TableName()},
		Body:  content,
	}

	searchResponse, err := search.Do(ctx, r.opensearchClient)
	if err != nil {
		return nil, err
	}
	if searchResponse.IsError() {
		return nil, fmt.Errorf("search error: %s", searchResponse.String())
	}

	defer searchResponse.Body.Close()

	dataB, err := io.ReadAll(searchResponse.Body)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Channel, 0)
	if err := json.Unmarshal(dataB, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes a channel from the channel repository.
//
// It takes a context.Context and the ID of the channel to be deleted as parameters.
// It returns an error if the deletion fails.
func (r *channelRepository) Delete(ctx context.Context, id int64) error {
	e := &entity.Channel{}
	del := opensearchapi.DeleteRequest{
		Index:      e.TableName(),
		DocumentID: fmt.Sprint(id),
	}

	deleteResponse, err := del.Do(ctx, r.opensearchClient)

	if err != nil {
		return err
	}

	if deleteResponse.IsError() {
		return fmt.Errorf("delete error: %s", deleteResponse.String())
	}

	return nil
}
