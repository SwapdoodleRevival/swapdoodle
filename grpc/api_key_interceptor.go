package grpc

import (
	"context"
	"errors"

	"github.com/silver-volt4/swapdoodle/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func apiKeyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if len(globals.GRPCApiKey) > 0 {
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			apiKeyHeader := md.Get("X-API-Key")

			if len(apiKeyHeader) == 0 || apiKeyHeader[0] != globals.GRPCApiKey {
				return nil, errors.New("missing or invalid API key")
			}
		} else {
			return nil, errors.New("API key verification failed")
		}
	}

	h, err := handler(ctx, req)

	return h, err
}
