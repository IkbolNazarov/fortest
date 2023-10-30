package main

import (
	"context"
	"fmt"
	"fortest/core"
	"fortest/db"
	"fortest/routes"
	"log"
	"net/http"
	"os"
	"time"

	my_handlers "fortest/handlers"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

func main() {

	l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("New looger error %v", err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(l)

	db.InitDB()

	var router = routes.Init(my_handlers.New(core.New()))

	zap.S().Infof("Starting Server on address http://%s:%s", "localhost", ":8080")
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"Accept", "Authorization", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	cors := handlers.CORS(credentials, methods, origins, headers)
	srv := &http.Server{
		Addr:         "8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      cors(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Block until we receive our signal.
	// <-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	zap.S().Info("shutting down")
	os.Exit(0)

}
