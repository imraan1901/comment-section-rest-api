package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "db"

type Database struct {
	Client *sqlx.DB
}

func NewDatabase(ctx context.Context) (*Database, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "NewDatabase", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return &Database{}, fmt.Errorf("could not connect to the database: %w", err)
	}

	return &Database{
		Client: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "Ping", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	return d.Client.DB.PingContext(ctx)
}
