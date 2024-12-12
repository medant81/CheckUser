package mssql

import (
	"CheckUser/internal/config"
	"CheckUser/internal/models"
	"context"
	"database/sql"
	"fmt"
	mssql "github.com/microsoft/go-mssqldb"
	"net/url"
)

type Storage struct {
	db *sql.DB
	sp StorageProcedure
}

type StorageProcedure struct {
	NameSP    string
	NameParam string
	TvpType   string
}

func New(dbConfig config.DBConfig, sp config.StorageProcedureConfig) (*Storage, error) {
	const op = "storage.mssql.New"

	query := url.Values{}
	query.Add("database", dbConfig.Database)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(dbConfig.User, dbConfig.Password),
		Host:     dbConfig.Host,
		Path:     dbConfig.Path, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", u.String())

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
		sp: StorageProcedure{
			NameSP:    sp.NameSP,
			NameParam: sp.NameParam,
			TvpType:   sp.TvpType,
		}}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) UsersCheck(ctx context.Context, userID []int64) (models.UsersResult, error) {
	const op = "storage.mssql.UsersCheck"

	type tableTvp struct {
		Id int64
	}

	tableTVPDate := make([]tableTvp, 0)

	for _, element := range userID {
		tableTVPDate = append(tableTVPDate, tableTvp{Id: element})
	}

	tvpType := mssql.TVP{
		TypeName: s.sp.TvpType,
		Value:    tableTVPDate,
	}

	err := s.db.PingContext(ctx)
	if err != nil {
		return models.UsersResult{}, fmt.Errorf("%s: %w", op, err)
	}

	tsql := fmt.Sprintf("EXEC %s @%s", s.sp.NameSP, s.sp.NameParam)

	stmt, err := s.db.Prepare(tsql)
	if err != nil {
		return models.UsersResult{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, sql.Named("telegram_id", tvpType))

	Users := make(map[int64]bool, 0)
	for rows.Next() {
		var telegram_id, is_check int64
		errScan := rows.Scan(&telegram_id, &is_check)
		if errScan == nil {
			Users[telegram_id] = (is_check == 1)
		}
	}

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
