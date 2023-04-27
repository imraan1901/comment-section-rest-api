package main

import (
	"fmt"

	"github.com/imraan1901/comment-section-rest-api/internal/comment/db"
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

	return nil

}

func main() {
	fmt.Println("GO REST API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
