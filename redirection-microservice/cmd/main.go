package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"redirection/internal/server"
	"redirection/internal/service"
	"redirection/proto/storagepb"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO later: RabbitMQ -> Analytics Microservice

func main() {
	log.Println("Loading config...")
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Loaded config successfully!")

	log.Println("Dialing storage service...")
	conn, err := grpc.NewClient(config.StorageAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to dial storage service:", err)
	}
	defer conn.Close()
	log.Println("Dialed storage service successfully!")

	cl := storagepb.NewStorageServiceClient(conn)

	s := service.NewRedirectionService(cl)

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
	ServerPort  string `envconfig:"PORT" required:"true"`
	StorageAddr string `envconfig:"STORAGE_ADDR" required:"true"`
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
