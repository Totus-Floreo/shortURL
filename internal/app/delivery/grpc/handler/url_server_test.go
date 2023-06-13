package grpchandler

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/Totus-Floreo/shortURL/internal/app/delivery/grpc/helpers"
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	"github.com/Totus-Floreo/shortURL/internal/app/domain/mocks"
	pb "github.com/Totus-Floreo/shortURL/internal/app/domain/proto"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context, service domain.IUrlService) (pb.ShortUrlClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterShortUrlServer(baseServer, NewShortUrlServer(service))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewShortUrlClient(conn)

	return client, closer
}

func TestCreateUrl(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockIUrlService(ctrl)

	client, closer := server(ctx, service)
	defer closer()

	type expectation struct {
		out          *pb.Short
		err          error
		serviceError error
	}

	tests := map[string]struct {
		in       *pb.Long
		expected expectation
	}{
		"Success": {
			in: &pb.Long{
				Link: "google.com",
			},
			expected: expectation{
				out: &pb.Short{
					Link: "bE2bqvWHr9",
				},
				err:          nil,
				serviceError: nil,
			},
		},
		"Error": {
			in: &pb.Long{
				Link: "google.com",
			},
			expected: expectation{
				out: &pb.Short{
					Link: "",
				},
				err:          helpers.GRPCError(domain.ErrorGenerateTimeout),
				serviceError: domain.ErrorGenerateTimeout,
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {

			service.EXPECT().CreateUrl(gomock.Any(), test.in.Link).Return(test.expected.out.Link, test.expected.serviceError)

			out, err := client.CreateUrl(ctx, test.in)
			if err != nil {
				if test.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", test.expected.err, err)
				}
			} else {
				if test.expected.out.Link != out.Link {
					t.Errorf("Out -> \nWant: %q\nGot : %q", test.expected.out, out)
				}
			}
		})
	}
}

func TestGetUrl(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockIUrlService(ctrl)

	client, closer := server(ctx, service)
	defer closer()

	type expectation struct {
		out          *pb.Long
		err          error
		serviceError error
	}

	tests := map[string]struct {
		in       *pb.Short
		expected expectation
	}{
		"Success": {
			in: &pb.Short{
				Link: "bE2bqvWHr9",
			},
			expected: expectation{
				out: &pb.Long{
					Link: "google.com",
				},
				err:          nil,
				serviceError: nil,
			},
		},
		"Error": {
			in: &pb.Short{
				Link: "bE2bqvWHr9",
			},
			expected: expectation{
				out: &pb.Long{
					Link: "",
				},
				err:          helpers.GRPCError(domain.ErrorGenerateTimeout),
				serviceError: domain.ErrorGenerateTimeout,
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {

			service.EXPECT().GetUrl(gomock.Any(), test.in.Link).Return(test.expected.out.Link, test.expected.serviceError)

			out, err := client.GetUrl(ctx, test.in)
			if err != nil {
				if test.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", test.expected.err, err)
				}
			} else {
				if test.expected.out.Link != out.Link {
					t.Errorf("Out -> \nWant: %q\nGot : %q", test.expected.out, out)
				}
			}
		})
	}
}
