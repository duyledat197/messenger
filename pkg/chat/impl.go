package chat

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"golang.org/x/sync/errgroup"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/websocket"
	"openmyth/messgener/util/snowflake"
)

// ChatImpl represents the implementation of the chat service.
type ChatImpl struct {
	onlineRepo  repository.OnlineRepository
	messageRepo repository.MessageRepository

	idGenerator *snowflake.Generator
}

// NewChatImpl creates a new instance of the ChatImpl struct, which implements the websocket.Impl interface for the Message type.
//
// It returns a pointer to the ChatImpl struct.
func NewChatImpl(
	onlineRepo repository.OnlineRepository,
	messageRepo repository.MessageRepository,
	idGenerator *snowflake.Generator,
) websocket.Impl[Message] {
	return &ChatImpl{
		onlineRepo:  onlineRepo,
		messageRepo: messageRepo,
		idGenerator: idGenerator,
	}
}

// Execute executes the chat message for the given client and data.
//
// It takes a client of type websocket.Client[Message] and a data of type Message as parameters.
// It returns an error.
func (i *ChatImpl) Execute(client *websocket.Client[Message], data Message) error {
	ctx := context.Background()
	toUserID := data.To

	online, err := i.onlineRepo.RetrieveByUserID(ctx, toUserID)
	if err != nil && !errors.Is(err, gocql.ErrNotFound) {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)
	id := i.idGenerator.Generate().Int64()

	eg.Go(func() error {
		return i.messageRepo.Create(ctx, &entity.Message{
			Bucket:    snowflake.MakeBucket(id),
			UserID:    client.UserID,
			FromID:    client.UserID,
			ToID:      toUserID,
			MessageID: id,
			Content:   data.Content,
			Reaction:  []string{},
		})
	})

	eg.Go(func() error {
		return i.messageRepo.Create(ctx, &entity.Message{
			Bucket:   snowflake.MakeBucket(id),
			UserID:   toUserID,
			FromID:   client.UserID,
			ToID:     toUserID,
			Content:  data.Content,
			Reaction: []string{},
		})
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if online != nil {
		client.Engine.Send <- websocket.SendInfo[Message]{
			ToClientID: online.ClientID,
			Data:       data,
		}
	}

	return nil
}

// Register registers a client with a user ID in the ChatImpl.
//
// Parameters: clientID string, userID string.
// Return type: error.
func (i *ChatImpl) Register(clientID, userID string) error {
	ctx := context.Background()
	if err := i.onlineRepo.Create(ctx, &entity.Online{
		UserID:   userID,
		ClientID: clientID,
	}); err != nil {
		return err
	}

	return nil
}

// UnRegister deletes the online status of a user.
//
// It takes a userID string as a parameter and returns an error if there was a problem deleting the online status.
func (i *ChatImpl) UnRegister(userID string) error {
	ctx := context.Background()
	if err := i.onlineRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	return nil
}
