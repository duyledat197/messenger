package service

import (
	"context"

	pb "openmyth/messgener/pb/chat"
)

type messageService struct {
	pb.UnimplementedMessageServiceServer
}

func NewMessageService() pb.MessageServiceServer {
	return &messageService{}
}

func (s *messageService) SendMessage(_ context.Context, _ *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *messageService) GetMessageLisyChannel(_ context.Context, _ *pb.GetMessageLisyChannelRequest) (*pb.GetMessageLisyChannelResponse, error) {
	panic("not implemented") // TODO: Implement
}
