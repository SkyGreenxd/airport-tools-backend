package postgres

import (
	"airport-tools-backend/pkg/e"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase struct {
	Db *gorm.DB
}

func Connect() (*PgDatabase, error) {
	const op = "Connect"
	path := os.Getenv("DB_URL")

	db, err := gorm.Open(pg.Open(path), &gorm.Config{})
	if err != nil {
		return nil, e.WrapWithFunc(op, "failed to connect database", err)
	}

	return &PgDatabase{Db: db}, nil
}

func (pg *PgDatabase) Ping() error {
	const op = "Ping"
	sqlDb, err := pg.Db.DB()
	if err != nil {
		return e.WrapWithFunc(op, "failed to get sql db instance", err)
	}

	if err := sqlDb.Ping(); err != nil {
		return e.WrapWithFunc(op, "failed to ping sql db", err)
	}

	return nil
}

func (pg *PgDatabase) Close() error {
	const op = "Close"

	sqlDb, err := pg.Db.DB()
	if err != nil {
		return e.WrapWithFunc(op, "failed to get sql db instance", err)
	}

	if err := sqlDb.Close(); err != nil {
		return e.WrapWithFunc(op, "failed to close sql db", err)
	}

	return nil
}

func (pg *PgDatabase) RunMigrations() error {
	const op = "RunMigrations"

	sqlDb, err := pg.Db.DB()
	if err != nil {
		return e.WrapWithFunc(op, "failed to get sql db instance", err)
	}

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return e.WrapWithFunc(op, "failed to create migrate driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)
	if err != nil {
		return e.WrapWithFunc(op, "failed to create migrate instance", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return e.WrapWithFunc(op, "migration failed", err)
	}

	log.Println("migrations applied successfully")
	return nil
}
