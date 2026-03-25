package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

const versionTable = "db_version"

type Migrator struct {
	migrator *migrate.Migrator
	pgxConn  *pgx.Conn
}

//go:embed data/*.sql
var migrationFiles embed.FS

func NewMigrator(dbDNS string) (Migrator, error) {
	conn, err := pgx.Connect(context.Background(), dbDNS)
	if err != nil {
		return Migrator{}, err
	}
	migrator, err := migrate.NewMigratorEx(
		context.Background(), conn, versionTable,
		&migrate.MigratorOptions{
			DisableTx: false,
		},
	)
	if err != nil {
		return Migrator{}, err
	}

	migrationRoot, err := fs.Sub(migrationFiles, "data")
	if err != nil {
		return Migrator{}, err
	}

	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		return Migrator{}, err
	}

	return Migrator{
		migrator: migrator,
		pgxConn:  conn,
	}, nil
}

func (m Migrator) Info() (int32, int32, string, error) {
	version, err := m.migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return 0, 0, "", err
	}
	info := ""

	var last int32
	for _, thisMigration := range m.migrator.Migrations {
		last = thisMigration.Sequence
		curr := version == thisMigration.Sequence

		indicator := "  "
		if curr {
			indicator = "->"
		}
		info = info + fmt.Sprintf(
			"%2s %3d %s\n",
			indicator,
			thisMigration.Sequence,
			thisMigration.Name,
		)
	}

	return version, last, info, nil
}

// Migrate the DB to the most latest migration
func (m Migrator) Migrate() error {
	return m.migrator.Migrate(context.Background())
}

func (m Migrator) ReleaseConn() error {
	return m.pgxConn.Close(context.Background())
}
