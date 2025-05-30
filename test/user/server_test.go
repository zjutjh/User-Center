package user

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	userv1 "github.com/zjutjh/User-Center-grpc/api/user/v1alpha1"
	serverv1 "github.com/zjutjh/User-Center-grpc/api/v1"
)

// 2. Test the gRPC SayHello method
func TestSayHello(t *testing.T) {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := serverv1.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Hello(ctx, &userv1.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.GetMessage())
	assert.Equal(t, "Hello World", r.GetMessage())
}
