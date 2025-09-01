package repository

import (
	"database/sql"
	"errors"
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
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	err = applyMigrations(db, logger)
	if err != nil {
		return nil, err
	}
	return &DBRepo{db: db}, nil
}

func (m *DBRepo) Put(url string) (model.URLID, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	done := false
	defer func() {
		if done {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	res, err := doPut(url, tx)

	done = true
	return res, err
}

type DB interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func doPut(url string, db DB) (model.URLID, error) {
	var urlFK int64
	var urlid model.URLID
	err := db.QueryRow("INSERT INTO urls (url) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", url).Scan(&urlFK)
	if err == nil {
		err = db.QueryRow("INSERT INTO links (url_fk) VALUES ($1) RETURNING id", urlFK).Scan(&urlid)
		if err != nil {
			return 0, fmt.Errorf("failed to insert link: %w", err)
		}
		return urlid, err
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to insert url: %w", err)
	}

	err = db.QueryRow("SELECT l.id FROM links l JOIN urls u ON l.url_fk = u.id WHERE u.url = $1", url).Scan(&urlid)
	if err != nil {
		return 0, fmt.Errorf("url is duplicated, but unable to get existing: %w", err)
	}

	return urlid, errs.NewDuplicatedURLError(url)
}

func (m *DBRepo) Get(id model.URLID) (string, error) {
	var url string
	err := m.db.QueryRow("SELECT u.url FROM links l JOIN urls u ON l.url_fk = u.id WHERE l.id = $1", id).Scan(&url)

	if err == sql.ErrNoRows {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	return url, err
}

func (m *DBRepo) BatchPut(urls []string) ([]model.URLID, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	done := false
	defer func() {
		if done {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	var res []model.URLID
	for _, url := range urls {
		id, err1 := doPut(url, tx)
		err = errors.Join(err, err1)
		res = append(res, id)
	}

	if err != nil {
		return nil, err
	}

	done = true
	return res, err
}

func (m *DBRepo) Ping() error {
	return m.db.Ping()
}

func applyMigrations(db *sql.DB, logger *zap.Logger) error {
	logger.Info("Applying migrations...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to init driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to init migrate: %w", err)
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
