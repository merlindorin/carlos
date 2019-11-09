package client

import (
	"github.com/iam-merlin/carlos"
	grpc2 "google.golang.org/grpc"
)

type Client struct {
	Car main.CarServiceClient
	Log main.LogServiceClient
}

func NewClient(conn *grpc2.ClientConn) Client {
	c := main.NewCarServiceClient(conn)
	l := main.NewLogServiceClient(conn)

	return Client{
		Car: c,
		Log: l,
	}
}
