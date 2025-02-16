package notification

import (
	"fmt"

	"github.com/leetcode-golang-classroom/golang-sample-with-grpc-server-streaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(gRPCAddress string) (proto.NotificationServiceClient, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(gRPCAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to construct client: %v", err)
	}

	return proto.NewNotificationServiceClient(conn), nil
}
