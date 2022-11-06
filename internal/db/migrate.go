package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func (d *DataBase) MigrateDB() error {

	driver, err := postgres.WithInstance(d.Client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create pg driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err.Error() == migrate.ErrNoChange.Error() {
		} else {
			return err
		}
	}
	log.Println("Successfully migrated database")
	return nil
}
