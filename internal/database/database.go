package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Pool  *pgxpool.Pool
	sqlDB *sql.DB
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Connect(connStr string) (*Database, error) {
	ctx := context.Background()
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB := stdlib.OpenDB(*poolConfig.ConnConfig)
	return &Database{Pool: pool, sqlDB: sqlDB}, nil
}

func (db *Database) RunMigrationsUp(dbName string) error {
	return db.runMigrations(dbName, true)
}

func (db *Database) RunMigrationsDown(dbName string) error {
	return db.runMigrations(dbName, false)
}

func (d *Database) runMigrations(dbName string, isUp bool) error {
	iofsDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}
	dbDriver, err := postgres.WithInstance(d.sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", iofsDriver, dbName, dbDriver)
	if err != nil {
		return fmt.Errorf("failed to migrate instance: %w", err)
	}
	if isUp {
		err = m.Up()
	} else {
		err = m.Down()
	}
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
