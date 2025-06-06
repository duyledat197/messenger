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

type channelService struct {
	channelRepo      repository.ChannelRepository
	cacheChannelRepo repository.CacheChannelRepository
	idGenerator      snowflake.Generator

	pb.UnimplementedChannelServiceServer
}

// NewChannelService creates a new instance of the pb.ChannelServiceServer interface.
// It takes a channelRepo of type repository.ChannelRepository and an idGenerator of type snowflake.Generator as parameters.
// It returns a pointer to a channelService struct that implements the pb.ChannelServiceServer interface.
func NewChannelService(
	channelRepo repository.ChannelRepository,
	cacheChannelRepo repository.CacheChannelRepository,
	idGenerator snowflake.Generator,
) pb.ChannelServiceServer {
	return &channelService{
		idGenerator:      idGenerator,
		channelRepo:      channelRepo,
		cacheChannelRepo: cacheChannelRepo,
	}
}

// SearchChannelByName searches for channels by name in the channel repository.
// It takes a context.Context and a SearchChannelByNameRequest as parameters.
// It returns a SearchChannelByNameResponse and an error.
func (s *channelService) SearchChannelByName(ctx context.Context, req *pb.SearchChannelByNameRequest) (*pb.SearchChannelByNameResponse, error) {
	channels, err := s.channelRepo.SearchByName(ctx, req.GetName(), req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to search channels: %v", err)
	}

	return &pb.SearchChannelByNameResponse{
		Channels: channelListToPbList(channels),
	}, nil
}

// channelListToPbList converts a list of entity.Channel pointers to a list of pb.Channel pointers.
func channelListToPbList(channels []*entity.Channel) []*pb.Channel {
	var result []*pb.Channel
	for _, c := range channels {
		result = append(result, channelToPb(c))
	}

	return result
}

// channelToPb converts an entity.Channel to a pb.Channel.
func channelToPb(channel *entity.Channel) *pb.Channel {
	return &pb.Channel{
		ChannelId:   channel.ChannelID,
		Name:        channel.Name,
		Description: channel.Description,
	}
}

// GetListChannel retrieves a list of channels from the channel repository.
// It takes a context.Context and a GetListChannelRequest as parameters.
// It returns a GetListChannelResponse and an error.
func (s *channelService) GetListChannel(ctx context.Context, req *pb.GetListChannelRequest) (*pb.GetListChannelResponse, error) {
	channels, err := s.cacheChannelRepo.List(ctx, req.Offset, req.Limit)
	if err == nil {
		return &pb.GetListChannelResponse{
			Channels: channelListToPbList(channels),
		}, nil
	}

	channels, err = s.channelRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to list channels: %v", err)
	}

	if err := s.cacheChannelRepo.CreateByList(ctx, req.Offset, req.Limit, channels); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to cache channels: %v", err)
	}

	return &pb.GetListChannelResponse{
		Channels: channelListToPbList(channels),
	}, nil
}
