package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rigidsearch/constants"
	"rigidsearch/indexing"
	"rigidsearch/router"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, skipping...")
	}
	log.Println("Rigid search starting up...")
	err := constants.LoadConstants()
	if err != nil {
		log.Fatalf("Failed to load constants: %v\n", err)
	}
	log.Println("Loaded constants")
	err = indexing.LoadIndex()
	if err != nil {
		log.Fatalf("Failed to load index: %v\n", err)
	}
	log.Println("Loaded index")
	router := router.NewRouter()
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	go func() {
		log.Println("Starting server on port 8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error on server listen: %v\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server shutdown err: ", err)
	}
	log.Println("Server exiting")
	err = indexing.StoreIndex()
	if err != nil {
		log.Println("Error storing index: ", err)
	}
	log.Println("Index stored")
	log.Println("Rigidsearch shut down")
}
