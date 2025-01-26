package auth

import (
	"context"
	ssov "github.com/MikeThesol/proto/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov.LoginRequest) (*ssov.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *ssov.RegisterRequest) (*ssov.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov.IsAdminRequest) (*ssov.IsAdminResponse, error) {
	panic("implement me")
}
