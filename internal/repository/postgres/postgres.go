package postgres

import (
	"airport-tools-backend/pkg/e"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDatabase struct {
	Db *gorm.DB
}

func Connect() (*PgDatabase, error) {
	const op = "Connect"
	path := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(path), &gorm.Config{})
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
