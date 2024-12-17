package app

import (
	grpcapp "CheckUser/internal/app/grpc"
	"CheckUser/internal/config"
	"CheckUser/internal/servises/check"
	"CheckUser/internal/storage/mssql"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cf *config.Config) *App {

	storage, err := mssql.New(cf.Storage, cf.StorageProcedure)

	if err != nil {
		panic(err)
	}

	checkService := check.New(log, storage, storage, cf.TokenTTL)

	grpcApp := grpcapp.New(log, checkService, cf.GRPC.Port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
