package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	pb "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	usercenter "github.com/zjutjh/User-Center/mygo_plugin"
	"github.com/zjutjh/User-Center/register"
	"github.com/zjutjh/mygo/foundation/kernel"
)

func TestMain(m *testing.M) {
	kernel.Bootstrap("../conf", register.Boot)
	go register.StartGrpcServer()
	m.Run()
}

func TestLogin(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	r, err := usercenter.Pick().Login(ctx, &pb.LoginRequest{StudentId: "teststudent", Password: "testpassword"})
	assert.NoError(t, err)
	t.Log(r)
}
