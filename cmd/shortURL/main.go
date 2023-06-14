package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	grpchandler "github.com/Totus-Floreo/shortURL/internal/app/delivery/grpc/handler"
	route "github.com/Totus-Floreo/shortURL/internal/app/delivery/http/handler"
	"github.com/Totus-Floreo/shortURL/internal/app/domain"
	pb "github.com/Totus-Floreo/shortURL/internal/app/domain/proto"
	"github.com/Totus-Floreo/shortURL/internal/app/repository/inmemory"
	"github.com/Totus-Floreo/shortURL/internal/app/repository/postgresql"
	"github.com/Totus-Floreo/shortURL/internal/app/service"
	"github.com/Totus-Floreo/shortURL/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("Starting...\n")
	zerologger := zerolog.New(os.Stderr)

	dbType := flag.String("dbType", "pgx", "Type of database (pgx or inmemory)")
	flag.Parse()

	var db domain.IUrlStorage
	switch *dbType {
	case "inmemory":
		db = inmemory.NewUrlStorage()
	case "pgx":
		postgreUrl := "postgres://" + os.Getenv("pg_url")
		pgxpool, err := pgxpool.New(context.Background(), postgreUrl)
		if err != nil {
			log.Fatalf("Postgre connection error: %v\n", err)
		}
		db = postgresql.NewUrlStorage(pgxpool)
	default:
		log.Fatalf("Unexpected dbType: %s\n", *dbType)
	}

	generator := service.NewGenerateLinkService()
	service := service.NewUrlService(db, generator)
	handlers := route.NewUrlHandler(service)
	grpcHandler := grpchandler.NewShortUrlServer(service)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", os.Getenv("gRPCport")))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(logger.InterceptorLogger(zerologger), opts...),
		),
	)
	pb.RegisterShortUrlServer(grpcServer, grpcHandler)

	go func() {
		log.Printf("Serve gRPC server on %s", os.Getenv("gRPCport"))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve grpc: %v", err)
		}
	}()

	router := gin.Default()
	router.GET("/:link", handlers.GetUrl)
	router.POST("/", handlers.CreateUrl)

	if err := router.Run(os.Getenv("httpport")); err != nil {
		log.Fatalf("Failed to serve http: %v", err)
	}
}
