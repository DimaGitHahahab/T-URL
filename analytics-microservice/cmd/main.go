package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"analytics/internal/repository"
	"analytics/internal/server"
	"analytics/internal/service"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	log.Println("Loading config...")
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Loaded config successfully!")

	repo := repository.NewRepository()

	s := service.New(repo)

	srv := server.New(s)

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("Server is running on port:", config.ServerPort)
		log.Fatal(srv.Run(config.ServerPort))
	}()

	<-sigQuit
	log.Println("Gracefully shutting down server...")
	srv.Stop()

	log.Println("Server shutdown is successful!")

}

type Config struct {
	ServerPort string `envconfig:"PORT" required:"true"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process env variables: %w", err)
	}

	return &cfg, nil
}
