package main

import (
	"context"
	"fmt"

	"github.com/imraan1901/comment-section-rest-api/internal/comment"
	"github.com/imraan1901/comment-section-rest-api/internal/db"
)

// Run - is responsible for
// the instantiation and startup of our
// go application
func Run() error {

	fmt.Println("Starting up our application")

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

	cmtService := comment.NewService(db)

	cmtService.PostComment(
		context.Background(),
		comment.Comment{
			ID: "f30b6b33-5351-4112-a20c-fc98c9319a73",
			Slug: "Manual test",
			Author: "Imraan",
			Body: "Hello world",
		},
	)

	fmt.Println(cmtService.GetComment(
		context.Background(),
		"f30b6b33-5351-4112-a20c-fc98c9319a73",
	))

	return nil

}

func main() {
	fmt.Println("GO REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
