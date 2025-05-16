package usersgrpc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/PratikKumar125/go-microservices/users/internal/service"
)

type UserGrpcServer struct {
	UserService *service.UserService
	UnimplementedUserRpcServiceServer
}

func NewUserGrpcServer(service *service.UserService) *UserGrpcServer {
	return &UserGrpcServer{
		UserService: service,
	}
}

func (userGrpcService *UserGrpcServer) GetUserByEmail(ctx context.Context, input *GetUserByEmailInput) (*GetUserByEmailResponse, error) {
	query := url.Values{}
	query.Set("email", input.Email)
	query.Set("name", input.Name)

	fmt.Printf("RPC REQUEST RECEIVED %v name %v email", input.Name, input.Email)

	users, err := userGrpcService.UserService.GetUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return &GetUserByEmailResponse{Users: []*User{}}, nil
	}

	rpcUsers := make([]*User, len(users))
	for i, user := range users {
		rpcUsers[i] = &User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return &GetUserByEmailResponse{Users: rpcUsers}, nil
}
