package auth

import (
	"context"
	ssov "github.com/MikeThesol/proto/proto/sso/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context, email string, password string, appID int32) (string, error)
	RegisterNewUser(
		ctx context.Context,
		name string,
		email string,
		password string,
	) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

const emptyValue = 0

type serverAPI struct {
	ssov.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov.LoginRequest) (*ssov.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov.LoginResponse{Token: token, Message: "Login is correct"}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov.RegisterRequest) (*ssov.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	_, err := s.auth.RegisterNewUser(ctx, req.GetName(), req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	token, err := s.auth.Login(ctx, req.GetName(), req.GetEmail(), 1) // TODO: изменить 1 - на нормальный id

	return &ssov.RegisterResponse{Token: token}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov.IsAdminRequest) (*ssov.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateIsAdmin(req *ssov.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "is_admin is required")
	}
	return nil
}

func validateRegister(req *ssov.RegisterRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateLogin(req *ssov.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}
