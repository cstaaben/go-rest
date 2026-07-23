/*
 * go-rest - A TUI for a REST client
 * Copyright (C) 2024  Corbin Staaben
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package db contains database logic.
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/cstaaben/go-rest/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	db *sqlx.DB
}

func New(ctx context.Context) (*DB, error) {
	db := &DB{}

	var (
		dbName string
		err    error

		dir        = config.DataDir()
		devEnabled = viper.GetBool("dev")
	)
	if devEnabled {
		dbName = "development"
	} else {
		dbName = "data"
	}

	db.db, err = sqlx.ConnectContext(ctx, "sqlite3", filepath.Join(dir, fmt.Sprintf("%s.db", dbName)))
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	if err := runMigrations(ctx, db.db.DB); err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}

	return db, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("Failed to set migration dialect", slog.String("error", err.Error()))
		return err
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		slog.Error("Failed to migrate database", slog.String("error", err.Error()))
		return err
	}

	return nil
}
