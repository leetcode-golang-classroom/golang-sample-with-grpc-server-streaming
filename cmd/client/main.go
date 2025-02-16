package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/internal/service/notification"
	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/proto"
)

func main() {
	// setup signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	// setup gRPC client and receive notification
	client, err := notification.NewClient(config.AppConfig.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	// ctx := context.Background()
	stream, err := client.GetNotifications(ctx, &proto.NotificationRequest{
		UserId: "123",
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		notification, err := stream.Recv()
		if err == io.EOF { // no more data to read
			break
		}
		if err != nil {
			log.Fatalf("failed to read notification: %v", err)
		}
		b, err := json.MarshalIndent(notification, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
