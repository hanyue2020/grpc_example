package service

import (
	"context"

	"example/api"
)

type Service struct {
	api.UnimplementedHelloServiceServer
}

func (s Service) Hello(ctx context.Context, request *api.HelloRequest) (*api.HelloResponse, error) {
	//TODO implement me
	println("Hello")
	return &api.HelloResponse{Data: request.Name}, nil
}

var _ api.HelloServiceServer = &Service{}
