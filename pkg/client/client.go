package client

import (
	"github.com/iam-merlin/carlos/grpc"
	grpc2 "google.golang.org/grpc"
)

type Client struct {
	Car grpc.CarServiceClient
	Log grpc.LogServiceClient
}

func NewClient(conn *grpc2.ClientConn) Client {
	c := grpc.NewCarServiceClient(conn)
	l := grpc.NewLogServiceClient(conn)

	return Client{
		Car: c,
		Log: l,
	}
}
