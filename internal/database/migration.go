package database

import (
	"echo-sqlc-template/internal/config"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"os"
)

const versionTable = "public.schema_version"

type Migrator struct {
	migrator *migrate.Migrator
}

func NewMigrator(dbDNS string) (Migrator, error) {
	conn, err := pgx.Connect(ctx, dbDNS)
	if err != nil {
		return Migrator{}, err
	}
	migrator, err := migrate.NewMigratorEx(
		ctx, conn, versionTable,
		&migrate.MigratorOptions{
			DisableTx: false,
		})
	if err != nil {
		return Migrator{}, err
	}

	migrationRoot := os.DirFS(config.Data.Database.Migrations)

	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		return Migrator{}, err
	}

	return Migrator{
		migrator: migrator,
	}, nil
}

func (m Migrator) Info() (int32, int32, string, error) {
	version, err := m.migrator.GetCurrentVersion(ctx)
	if err != nil {
		return 0, 0, "", err
	}
	info := ""

	var last int32
	for _, thisMigration := range m.migrator.Migrations {
		last = thisMigration.Sequence

		cur := version == thisMigration.Sequence
		indicator := "  "
		if cur {
			indicator = "->"
		}
		info = info + fmt.Sprintf(
			"%2s %3d %s\n",
			indicator,
			thisMigration.Sequence, thisMigration.Name)
	}

	return version, last, info, nil
}

func (m Migrator) Migrate() error {
	err := m.migrator.Migrate(ctx)
	return err
}
