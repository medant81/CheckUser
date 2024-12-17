package check

import (
	"CheckUser/internal/lib/jwt"
	"CheckUser/internal/lib/logger/sl"
	"CheckUser/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Check struct {
	log          *slog.Logger
	usersChecker UsersCheckerDB
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UsersCheckerDB interface {
	UsersCheckDB(ctx context.Context, userID []int64) (models.UsersResult, error)
	//TokenGet(ctx context.Context, userName string, password string, appID int) (string, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns a new instance of the Check service
func New(log *slog.Logger, usersChecker UsersCheckerDB, appProvider AppProvider, tokenTTL time.Duration) *Check {
	return &Check{
		log:          log,
		usersChecker: usersChecker,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (c *Check) TokenServises(ctx context.Context, userName string, password string, appID int) (string, error) {
	const op = "check.Token"

	log := c.log.With(
		slog.String("op", op),
		slog.String("username", userName),
	)

	log.Info("attempting to token")

	app, err := c.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(userName, password, app, c.tokenTTL)
	if err != nil {
		c.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (c *Check) CheckUsersServises(ctx context.Context, tokenRes string, userID []int64) (models.UsersResult, error) {
	const op = "check.UsersCheck"

	log := c.log.With("op", op)
	log.Info("Check users")

	_usersResult, err := c.usersChecker.UsersCheckDB(ctx, userID)
	if err != nil {
		return models.UsersResult{}, fmt.Errorf("%s: %w", op, err)
	}

	return _usersResult, nil
}
