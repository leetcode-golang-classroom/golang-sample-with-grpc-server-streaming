package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/internal/service/notification"
	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {

	// setup signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// construct gRPC service
	listener, err := net.Listen("tcp", config.AppConfig.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// setup redis client
	redisOpts, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(redisOpts)
	handler := notification.NewHandler(redisClient)
	proto.RegisterNotificationServiceServer(grpcServer, handler)
	slog.Info("listening on " + config.AppConfig.GRPCAddress)
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// listen for stop signal
	<-ctx.Done()
	_, stop := context.WithTimeout(ctx, time.Second*10)
	slog.Info("stopping server, wait for 10 seconds to stop")
	grpcServer.GracefulStop()
	defer stop()
}
