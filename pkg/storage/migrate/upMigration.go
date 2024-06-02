package migrate

import (
	"Ozon_testtask/internal/migrations"
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"log"
)

func UpMigration(ctx context.Context, db *sql.DB) error {
	provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)
	if err != nil {
		log.Fatal("Main failed to create NewProvider for migration")

		return err
	}

	_, err = provider.Up(ctx)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)

		return err
	}

	return nil
}
