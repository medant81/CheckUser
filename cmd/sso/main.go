package main

import (
	"CheckUser/internal/app"
	"CheckUser/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	test := test.Test{}
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting app: ",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPC.Port),
	)
	/*portStr := ""
	if cfg.Storage.Port > 0 {
		portStr = fmt.Sprintf(";port=%d", cfg.Storage.Port)
	}
	connString := fmt.Sprintf("server=%s%s;user id=%s;password=%s;database=%s;",
		cfg.Storage.Server,
		portStr,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Database)*/

	application := app.New(log, cfg.GRPC.Port, cfg.Storage, cfg.StorageProcedure, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCSrv.Stop()

	log.Info("gRPC app stop")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
