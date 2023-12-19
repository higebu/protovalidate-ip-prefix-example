package server

import (
	"context"

	examplev1 "github.com/higebu/protovalidate-ip-prefix-example/gen/go/proto/example/v1"
)

// Server implements ExampleService.
type Server struct {
	examplev1.UnimplementedExampleServiceServer
}

func New() *Server {
	return &Server{}
}

func (*Server) AddIPv4Address(ctx context.Context, req *examplev1.AddIPv4AddressRequest) (*examplev1.AddIPv4AddressResponse, error) {
	return &examplev1.AddIPv4AddressResponse{
		Name:    req.GetName(),
		Address: req.GetAddress(),
	}, nil
}

func (*Server) AddIPv6Address(ctx context.Context, req *examplev1.AddIPv6AddressRequest) (*examplev1.AddIPv6AddressResponse, error) {
	return &examplev1.AddIPv6AddressResponse{
		Name:    req.GetName(),
		Address: req.GetAddress(),
	}, nil
}

func (*Server) AddIPv4Route(ctx context.Context, req *examplev1.AddIPv4RouteRequest) (*examplev1.AddIPv4RouteResponse, error) {
	return &examplev1.AddIPv4RouteResponse{
		Name:   req.GetName(),
		Prefix: req.GetPrefix(),
	}, nil
}

func (*Server) AddIPv6Route(ctx context.Context, req *examplev1.AddIPv6RouteRequest) (*examplev1.AddIPv6RouteResponse, error) {
	return &examplev1.AddIPv6RouteResponse{
		Name:   req.GetName(),
		Prefix: req.GetPrefix(),
	}, nil
}
