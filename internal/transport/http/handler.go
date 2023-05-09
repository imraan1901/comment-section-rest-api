package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"

	tr "go.opentelemetry.io/otel/trace"
)

type Handler struct {
	Router  *mux.Router
	Service CommentService
	Server  *http.Server
}

// name is the Tracer name used to identify this instrumentation library.
const name = "http"

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive")
	})

	h.Router.HandleFunc("/api/v1/comment", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/v1/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/v1/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")

}

func (h *Handler) Serve(ctx context.Context) error {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "Serve", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	// This will run the go func in a thread ()
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	// Make a channel and wait here until an interrupt happens on this thread
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Once an interrupt happens shutdown the server in 15 seconds
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	// defer func gets called at the end of the 15 seconds to kill the main thread
	defer cancel()
	// Withing the 15 seconds process the remaining in flight requests
	h.Server.Shutdown(ctx)

	log.Println("Shutting down gracefully")

	return nil
}
