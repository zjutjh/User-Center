package user

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	userv1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
)

// 2. Test the gRPC SayHello method
func TestLogin(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := userv1.NewUserCenterServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Login(ctx, &userv1.LoginRequest{StudentId: "teststudent", Password: "testpassword"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.GetMessage())
}
