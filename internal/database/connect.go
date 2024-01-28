package database

import (
	"context"
	"echo-sqlc-template/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Q *Queries

var ctx = context.Background()

func Connect() error {
	conn, err := pgxpool.New(ctx, config.Data.Database.Uri)
	if err != nil {
		return err
	}

	migrator, err := NewMigrator(config.Data.Database.Uri)
	if err != nil {
		return err
	}

	err = migrator.Migrate()
	if err != nil {
		return err
	}

	Q = New(conn)

	return nil
}
