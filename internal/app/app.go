package app

import (
	grpcapp "CheckUser/internal/app/grpc"
	"CheckUser/internal/config"
	"CheckUser/internal/servises/check"
	"CheckUser/internal/storage/mssql"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, dbConfig config.DBConfig, sp config.StorageProcedureConfig, tokenTTL time.Duration) *App {

	storage, err := mssql.New(dbConfig, sp)
	//storage, err := mssql.New(cfg.Sto)
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
