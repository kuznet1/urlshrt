package repository

import (
	"context"
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

func (m *DBRepo) Put(ctx context.Context, url string) (model.URLID, error) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return 0, err
	}

	tx, err := m.db.BeginTx(ctx, nil)
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

	res, err := doPut(url, userId, tx)

	done = true
	return res, err
}

func doPut(url string, userID int, tx *sql.Tx) (model.URLID, error) {
	var urlid model.URLID
	err := tx.QueryRow("INSERT INTO links (url, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id", url, userID).Scan(&urlid)
	if err == nil {
		return urlid, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to insert url: %w", err)
	}

	err = tx.QueryRow("SELECT id FROM links WHERE url = $1", url).Scan(&urlid)
	if err != nil {
		return 0, fmt.Errorf("url is duplicated, but unable to get existing: %w", err)
	}

	return urlid, errs.NewDuplicatedURLError(url)
}

func (m *DBRepo) Get(ctx context.Context, id model.URLID) (string, error) {
	var url string
	err := m.db.QueryRowContext(ctx, "SELECT url FROM links WHERE id = $1", id).Scan(&url)

	if err == sql.ErrNoRows {
		return "", errs.NewHTTPError(fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
	}

	return url, err
}

func (m *DBRepo) BatchPut(ctx context.Context, urls []string) ([]model.URLID, error) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := m.db.BeginTx(ctx, nil)
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
		id, err1 := doPut(url, userId, tx)
		err = errors.Join(err, err1)
		res = append(res, id)
	}

	if err != nil {
		return nil, err
	}

	done = true
	return res, err
}

func (m *DBRepo) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

func (m *DBRepo) UserUrls(ctx context.Context) (map[model.URLID]string, error) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryContext(ctx, "SELECT id, url FROM links WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query urls: %w", err)
	}
	defer rows.Close()

	res := make(map[model.URLID]string)
	for rows.Next() {
		var id model.URLID
		var url string
		err = rows.Scan(&id, &url)
		if err != nil {
			return nil, fmt.Errorf("failed to scan url: %w", err)
		}
		res[id] = url
	}

	return res, nil
}

func (m *DBRepo) CreateUser(ctx context.Context) (int, error) {
	var userID int
	err := m.db.QueryRowContext(ctx, "INSERT INTO users DEFAULT VALUES RETURNING id").Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return userID, nil
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
