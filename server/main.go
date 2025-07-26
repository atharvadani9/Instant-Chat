package main

import (
	"chat/internal/app"
	"chat/routes"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go server port")
	flag.Parse()

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()

	r := routes.SetupRoutes(application)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	application.Logger.Printf("Starting server on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		application.Logger.Fatalf("ERROR: Server error %v", err)
	}
}
