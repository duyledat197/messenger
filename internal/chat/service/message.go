package service

import (
	"context"

	pb "openmyth/messgener/pb/chat"
)

type messageService struct {
	pb.UnimplementedMessageServiceServer
}

// NewMessageService returns a new instance of the pb.MessageServiceServer interface.
//
// It takes no parameters and returns a pointer to a messageService struct that implements the pb.MessageServiceServer interface.
func NewMessageService() pb.MessageServiceServer {
	return &messageService{}
}

func (s *messageService) SendMessage(_ context.Context, _ *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *messageService) GetMessageLisyChannel(_ context.Context, _ *pb.GetMessageLisyChannelRequest) (*pb.GetMessageLisyChannelResponse, error) {
	panic("not implemented") // TODO: Implement
}
