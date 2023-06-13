package helpers

import (
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCError(err error) error {
	switch err {
	case domain.ErrorGenerateTimeout:
		return status.Errorf(codes.Canceled, "Generate link timeout")
	case domain.ErrorLinkNotFound:
		return status.Errorf(codes.NotFound, "Link not found")
	case domain.ErrorInvalidLink:
		return status.Errorf(codes.InvalidArgument, "Link is broken")
	default:
		return status.Errorf(codes.Internal, "Internal Error")
	}
}
