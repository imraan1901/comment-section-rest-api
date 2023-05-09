package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

// name is the Tracer name used to identify this instrumentation library.

func (d *Database) MigrateDB(ctx context.Context) error {

	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	fmt.Println("migrating our database")

	driver, err := postgres.WithInstance(d.Client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	fmt.Println("successfully migrated the database")
	return nil
}
