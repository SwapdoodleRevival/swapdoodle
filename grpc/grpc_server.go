package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/silver-volt4/swapdoodle/globals"
	"google.golang.org/grpc"
)

type gRPCSwapdoodleServer struct {
	UnimplementedSwapdoodleServer
}

func StartGRPCServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", globals.GRPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(apiKeyInterceptor),
	)

	RegisterSwapdoodleServer(server, &gRPCSwapdoodleServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
