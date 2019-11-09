package client

import (
	"github.com/iam-merlin/carlos/internal/grpc"
	grpc2 "google.golang.org/grpc"
)

type Client struct {
	car grpc.CarServiceClient
	log grpc.LogServiceClient
}

func NewClient(conn *grpc2.ClientConn) Client {
	c := grpc.NewCarServiceClient(conn)
	l := grpc.NewLogServiceClient(conn)

	return Client{
		car: c,
		log: l,
	}
}
