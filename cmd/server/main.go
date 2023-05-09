package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/imraan1901/comment-section-rest-api/internal/comment"
	"github.com/imraan1901/comment-section-rest-api/internal/db"
	transportHttp "github.com/imraan1901/comment-section-rest-api/internal/transport/http"


	tr "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("main"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "main"),
		),
	)
	return r
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// name is the Tracer name used to identify this instrumentation library.
const name = "main"

// Run - is responsible for
// the instantiation and startup of our
// go application
func Run() error {

	l := log.New(os.Stdout, "", 0)
	// Added time format used by API requests
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	// Write telemetry data to a file.
	tracerFile := path.Join("tracers", now+"_traces.txt")
	f, err := os.Create(tracerFile)
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	startTime := time.Now()
	ctx, span := otel.Tracer(name).Start(context.Background(), "main", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	fmt.Println("Starting up our application")

	// DB layer
	db, err := db.NewDatabase(ctx)
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(ctx); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	fmt.Println("successfully connected and pinged database")

	// DB layer passed into business layer
	cmtService := comment.NewService(db)

	// business layer passed into transport/http layer
	httpHandler := transportHttp.NewHandler(cmtService)
	if err := httpHandler.Serve(ctx); err != nil {
		return err
	}

	return nil

}

func main() {
	fmt.Println("GO REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
