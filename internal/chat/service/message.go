package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openmyth/messgener/internal/chat/entity"
	"openmyth/messgener/internal/chat/repository"
	pb "openmyth/messgener/pb/chat"
	"openmyth/messgener/util/snowflake"
)

type messageService struct {
	messageRepo repository.MessageRepository
	channelRepo repository.ChannelRepository
	idGenerator snowflake.Generator

	pb.UnimplementedMessageServiceServer
}

// NewMessageService returns a new instance of the pb.MessageServiceServer interface.
// It takes no parameters and returns a pointer to a messageService struct that implements the pb.MessageServiceServer interface.
func NewMessageService(idGenerator snowflake.Generator, messageRepo repository.MessageRepository, channelRepo repository.ChannelRepository) pb.MessageServiceServer {
	return &messageService{
		messageRepo: messageRepo,
		channelRepo: channelRepo,
		idGenerator: idGenerator,
	}
}

// SendMessage sends a message and returns the message ID.
// It takes a context.Context and a *pb.SendMessageRequest as parameters.
// It returns a *pb.SendMessageResponse and an error.
func (s *messageService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	messageID := s.idGenerator.Generate().Int64()
	if err := s.messageRepo.Create(ctx, &entity.Message{
		ChannelID: req.ChannelId,
		Content:   req.Content,
		UserID:    req.UserId,
		MessageID: messageID,
		Bucket:    snowflake.MakeBucket(messageID),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	return &pb.SendMessageResponse{
		MessageId: messageID,
	}, nil
}

// GetMessageListChannel retrieves a list of messages for a given channel ID, starting from the specified offset and limited to the specified limit.
// It takes a context.Context and a *pb.GetMessageListChannelRequest as parameters.
// It returns a *pb.GetMessageListChannelResponse and an error.
func (s *messageService) GetMessageListChannel(ctx context.Context, req *pb.GetMessageListChannelRequest) (*pb.GetMessageListChannelResponse, error) {
	if _, err := s.channelRepo.RetrieveByChannelID(ctx, req.GetChannelId()); err != nil {
		return nil, status.Errorf(codes.NotFound, "channel not found")
	}

	result, err := s.messageRepo.RetrieveMessages(ctx, req.ChannelId, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve messages: %v", err)
	}

	return &pb.GetMessageListChannelResponse{
		Messages: messageListToPbList(result),
	}, nil
}

// messageListToPbList converts a list of entity.Message objects to a list of pb.Message objects.
func messageListToPbList(result []*entity.Message) []*pb.Message {
	var resp []*pb.Message
	for _, v := range result {
		resp = append(resp, messageToPb(v))
	}

	return resp
}

// messageToPb converts an entity.Message to a pb.Message.
func messageToPb(message *entity.Message) *pb.Message {
	return &pb.Message{
		UserId:    message.UserID,
		Content:   message.Content,
		MessageId: message.MessageID,
		Reaction:  message.Reaction,
	}
}
