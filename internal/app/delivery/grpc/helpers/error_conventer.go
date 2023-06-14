package helpers

import (
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCError(err error) error {
	switch err {
	case domain.ErrorGenerateTimeout:
		return status.Errorf(codes.Canceled, "Error: %s", err.Error())
	case domain.ErrorLinkNotFound:
		return status.Errorf(codes.NotFound, "Error: %s", err.Error())
	case domain.ErrorInvalidShort:
		return status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	case domain.ErrorInvalidLink:
		return status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	default:
		return status.Errorf(codes.Internal, "Error: %s", err.Error())
	}
}
