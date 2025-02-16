package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/internal/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	// setup signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	// publish data to the notification channel via redis
	opts, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(opts)
	channelName := fmt.Sprintf("notifications/%s", "123")
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			redisClient.Close()
			log.Println("shutting down")
			return
		case t := <-ticker.C:
			dataToSend := fmt.Sprintf("New notication %s", t.String())
			slog.Info("starting sending data: " + dataToSend)
			if cmd := redisClient.Publish(ctx, channelName, dataToSend); cmd.Err() != nil {
				log.Fatal(cmd.Err())
			}

		}
	}

}
