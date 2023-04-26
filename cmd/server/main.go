package main

import "fmt"

// Run - is responsible for 
// the instantiation and startup of our
// go application 
func Run() error {

	fmt.Println("Starting up our application")
	return nil

}

func main() {
	fmt.Println("API GO REST")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}