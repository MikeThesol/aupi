package app

import (
	grpcapp "auth-api/internal/app/grpc"
	"auth-api/internal/services/auth"
	"github.com/MikeThesol/proto/proto/user/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	userServiceAddr string,
	tokenTTL time.Duration,
) *App {
	userCon, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to user service", slog.String("error", err.Error()))
		panic(err)
	}

	userClient := gen.NewUserClient(userCon)

	authService := auth.New(log, userClient, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, tokenTTL, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
