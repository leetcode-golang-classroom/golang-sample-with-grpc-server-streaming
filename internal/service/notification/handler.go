package notification

import (
	"fmt"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var _ proto.NotificationServiceServer = (*Handler)(nil)

type Handler struct {
	proto.UnimplementedNotificationServiceServer
	redisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{
		redisClient: redisClient,
	}
}

func (h *Handler) GetNotifications(request *proto.NotificationRequest, stream grpc.ServerStreamingServer[proto.Notification]) error {
	pubsub := h.redisClient.Subscribe(stream.Context(),
		fmt.Sprintf("notifications/%s", request.GetUserId()))
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case msg := <-pubsub.Channel():

			if err := stream.Send(&proto.Notification{
				UserId: request.UserId,
				Content: fmt.Sprintf("New notification at %s: %s", time.Now().String(),
					msg.Payload,
				),
				CreatedAt: time.Now().UnixMilli(),
			}); err != nil {
				return fmt.Errorf("could not sent notification: %w", err)
			}
		}
	}
}
