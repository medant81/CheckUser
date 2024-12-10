package app

import (
	grpcapp "CheckUser/internal/app/grpc"
	"CheckUser/internal/servises/check"
	"CheckUser/internal/storage/mssql"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: init storage
	storage, err := mssql.New(storagePath)
	if err != nil {
		panic(err)
	}

	// TODO: init check service
	checkService := check.New(log, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, checkService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
