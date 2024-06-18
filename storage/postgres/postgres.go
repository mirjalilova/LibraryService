package postgres

import (
	"database/sql"
	"fmt"

	"golang.org/x/exp/slog"

	"library_service/config"
	"library_service/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	AuthorS   storage.AuthorI
	BookS     storage.BookI
	BorrowerS storage.BorrowerI
	GenreS    storage.GenreI
}

func Connect(cnf config.Config) (*Storage, error) {
	dbCon := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cnf.DB_USER, cnf.DB_PASSWORD, cnf.DB_HOST, cnf.DB_PORT, cnf.DB_NAME)
	db, err := sql.Open("postgres", dbCon)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("can't ping db: %v", err)
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	slog.Info("connected to db")

	return &Storage{
		db:        db,
		AuthorS:   NewAuthorRepo(db),
		BookS:     NewBookRepo(db),
		BorrowerS: NewBorrowerRepo(db),
		GenreS:    NewGenreRepo(db),
	}, nil
}
