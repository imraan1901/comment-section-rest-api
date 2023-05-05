package main

import (
	"fmt"

	"github.com/imraan1901/comment-section-rest-api/internal/comment"
	transportHttp "github.com/imraan1901/comment-section-rest-api/internal/transport/http"
	"github.com/imraan1901/comment-section-rest-api/internal/db"
)

// Run - is responsible for
// the instantiation and startup of our
// go application
func Run() error {

	fmt.Println("Starting up our application")

	// DB layer
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	fmt.Println("successfully connected and pinged database")

	// DB layer passed into business layer
	cmtService := comment.NewService(db)

	// business layer passed into transport/http layer
	httpHandler := transportHttp.NewHandler(cmtService)
	if err := httpHandler.Serve(); err != nil {
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
