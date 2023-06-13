package grpchandler

import (
	"context"

	"github.com/Totus-Floreo/shortURL/internal/app/delivery/grpc/helpers"
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	pb "github.com/Totus-Floreo/shortURL/internal/app/domain/proto"
)

type ShortUrlhandler struct {
	pb.UnimplementedShortUrlServer

	service domain.IUrlService
}

func NewShortUrlServer(service domain.IUrlService) *ShortUrlhandler {
	return &ShortUrlhandler{
		service: service,
	}
}

func (s *ShortUrlhandler) CreateUrl(ctx context.Context, long *pb.Long) (*pb.Short, error) {
	longUrl := long.GetLink()

	short, err := s.service.CreateUrl(ctx, longUrl)
	if err != nil {
		return nil, helpers.GRPCError(err)
	}

	return &pb.Short{Link: short}, nil
}

func (s *ShortUrlhandler) GetUrl(ctx context.Context, short *pb.Short) (*pb.Long, error) {
	shortUrl := short.GetLink()

	long, err := s.service.GetUrl(ctx, shortUrl)
	if err != nil {
		return nil, helpers.GRPCError(err)
	}

	return &pb.Long{Link: long}, nil
}
