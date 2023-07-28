package rpc

import "context"

type ArticleRPC struct{}

func (server *ArticleRPC) Hello(_ context.Context, _ *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	return &greeter.HelloResponse{
		Message: "hello!",
	}, nil
}
