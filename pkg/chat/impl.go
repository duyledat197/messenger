package chat

import (
	"context"
	"errors"

	"github.com/gocql/gocql"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	"openmyth/messgener/pkg/websocket"
	"openmyth/messgener/util/snowflake"
)

// Impl represents the implementation of the chat service.
type Impl struct {
	onlineRepo  repository.OnlineRepository
	messageRepo repository.MessageRepository
}

// NewImpl creates a new instance of the Impl struct, which implements the websocket.Impl interface for the Message type.
//
// It returns a pointer to the Impl struct.
func NewImpl(
	onlineRepo repository.OnlineRepository,
	messageRepo repository.MessageRepository,
) websocket.Impl[Message] {
	return &Impl{
		onlineRepo:  onlineRepo,
		messageRepo: messageRepo,
	}
}

// Execute executes the chat message for the given client and data.
//
// It takes a client of type websocket.Client[Message] and a data of type Message as parameters.
// It returns an error.
func (i *Impl) Execute(client *websocket.Client[Message], data Message) error {
	ctx := context.Background()
	toUserID := data.To

	online, err := i.onlineRepo.RetrieveByUserID(ctx, toUserID)
	if err != nil && !errors.Is(err, gocql.ErrNotFound) {
		return err
	}

	id := snowflake.GenerateID()

	if err := i.messageRepo.Create(ctx, &entity.Message{
		Bucket:    snowflake.MakeBucket(id),
		UserID:    client.UserID,
		ChannelID: client.ChannelID,
		MessageID: id,
		Content:   data.Content,
		Reaction:  []string{},
	}); err != nil {
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

// Register registers a client with a user ID in the Impl.
//
// Parameters: clientID string, userID string.
// Return type: error.
func (i *Impl) Register(clientID, userID string) error {
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
func (i *Impl) UnRegister(userID string) error {
	ctx := context.Background()
	if err := i.onlineRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	return nil
}
