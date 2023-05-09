package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

// name is the Tracer name used to identify this instrumentation library.

func (d *Database) MigrateDB(ctx context.Context) error {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "MigrateDB", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	fmt.Println("migrating our database")

	driver, err := postgres.WithInstance(d.Client.DB, &postgres.Config{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		fmt.Println(err)
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	fmt.Println("successfully migrated the database")
	return nil
}
