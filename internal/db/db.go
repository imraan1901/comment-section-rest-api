package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "db"


type Database struct {
	Client *sqlx.DB
}

func NewDatabase(ctx context.Context) (*Database, error) {

	_, span := otel.Tracer(name).Start(ctx, "db")
	defer span.End()

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
		return &Database{}, fmt.Errorf("could not connect to the database: %w", err)
	}

	return &Database{
		Client: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	_, span := otel.Tracer(name).Start(ctx, "ping")
	defer span.End()
	
	return d.Client.DB.PingContext(ctx)
}
