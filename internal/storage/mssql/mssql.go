package mssql

import (
	"CheckUser/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/microsoft/go-mssqldb/azuread"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.mssql.New"

	db, err := sql.Open(azuread.DriverName, storagePath)
	//db, err := sql.Open(azuread.DriverName, "server=TRAMPAMPAM\\SQL2017;user id=admingo;password=admingo;database=integ01;")
	//db, err := sql.Open(azuread.DriverName, "TRAMPAMPAM\\SQL2017;port=1433;database=integ01;fedauth=ActiveDirectoryDefault;")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) UsersCheck(ctx context.Context, userID []int64) (models.UsersResult, error) {
	const op = "storage.mssql.UsersCheck"

	err := s.db.PingContext(ctx)
	if err != nil {
		return models.UsersResult{}, fmt.Errorf("%s: %w", op, err)
	}

	tsql := `DECLARE	@return_value int
				EXEC	@return_value = [dbo].[BOT_AdminChat_Dismissed@check]
				@telegram_id = 132456789`

	stmt, err := s.db.Prepare(tsql)
	if err != nil {
		return models.UsersResult{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, sql.Named("telegram_id", 123465))
	//row := stmt.QueryRowContext(ctx, sql.Named("telegram_id", 132456789))

	//var _usersRes models.UsersResult
	Users := make(map[int64]bool, 0)
	for rows.Next() {
		var telegram_id, is_check int64
		err_scan := rows.Scan(&telegram_id, &is_check)
		if err_scan == nil {
			Users[telegram_id] = (is_check == 1)
		}
	}

	//var output any
	//err = row.Scan(output)

	return models.UsersResult{Users: Users}, nil
}

func (s *Storage) TokenGet(ctx context.Context, userName string, password string, appID int) (string, error) {
	return "", nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	// TODO: Доделать
	modelsApp := models.App{
		ID:     1,
		Name:   "test",
		Secret: "test",
	}

	return modelsApp, nil
}
