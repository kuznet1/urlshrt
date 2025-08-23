package repository

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"go.uber.org/zap"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBRepo struct {
	db *sql.DB
}

func NewDBRepo(dsn string, logger *zap.Logger) (*DBRepo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = applyMigrations(db, logger)
	if err != nil {
		return nil, err
	}
	return &DBRepo{db: db}, nil
}

func (m *DBRepo) Put(url string) (model.URLID, error) {
	var id model.URLID
	err := m.db.QueryRow("INSERT INTO links (url) VALUES ($1) RETURNING id", url).Scan(&id)
	return id, err
}

func (m *DBRepo) Get(id model.URLID) (string, error) {
	var url string
	err := m.db.QueryRow("SELECT url FROM links WHERE id = $1", id).Scan(&url)

	if err == sql.ErrNoRows {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	return url, err
}

func (m *DBRepo) Ping() error {
	return m.db.Ping()
}

func applyMigrations(db *sql.DB, logger *zap.Logger) error {
	logger.Info("Applying migrations...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	switch err {
	case nil:
		logger.Info("Migrations applied successfully.")
		return nil
	case migrate.ErrNoChange:
		logger.Info("Database is up to date.")
		return nil
	default:
		return fmt.Errorf("migration failed: %v", err)
	}
}
