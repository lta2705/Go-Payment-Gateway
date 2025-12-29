package main

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/app"
)

func main() {
	app, err := app.InitializeApp()
	if err != nil {
		panic(err)
	}

	// Chạy Web Server
	app.Router.Run(":8085")

}
