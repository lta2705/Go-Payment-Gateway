package main

import (
	"context"
	"github.com/lta2705/Go-Payment-Gateway/internal/app"
)

func main() {
	app, err := app.InitializeApp()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//initialize consumer worker on a separate goroutine
	go app.Consumer.ReadTransaction(ctx)

	// Chạy Web Server
	err = app.Router.Run(":8085")
	if err != nil {
		panic(err)
	}

}
