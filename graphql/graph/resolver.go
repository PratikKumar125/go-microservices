package graph

import (
	"context"
	"log"

	usersgrpc "github.com/PratikKumar125/go-microservices/users/usersrpc"
	"github.com/knadh/koanf/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Resolver struct {
	UserRpcClient usersgrpc.UserRpcServiceClient
}

func NewResolver(ctx context.Context) (*Resolver, error) {
	var k = koanf.New(".")
	appConfig, err := NewAppConfig(k, "/Users/pratikkumar/Downloads/Personal/golang/Microservices/graphql/graph/dev.env.yaml")
	if err != nil {
		panic(err)
	}

	// 1. Establish Users mircroservice gRPC connection
	grpcAddr := appConfig.ConfigService.String("app.users_grpc_addr")
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server at %s: %v", grpcAddr, err)
		return nil, err
	}

	// Create Users gRPC client
	client := usersgrpc.NewUserRpcServiceClient(conn)

	return &Resolver{
		UserRpcClient: client,
	}, nil
}
