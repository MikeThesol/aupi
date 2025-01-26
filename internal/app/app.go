package app

import (
	grpcapp "auth-api/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcPort, tokenTTL)
	return &App{
		GRPCSrv: grpcApp,
	}
}
