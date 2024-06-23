package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/handler"
	"api-gateway/internal/handlergen"
	"api-gateway/internal/service"
	"api-gateway/proto/redirectionpb"
	"api-gateway/proto/shorteningpb"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Loading config...")
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Loaded config successfully!")

	log.Println("Dialing shortening service...")
	shorteningConn, err := grpc.NewClient(config.ShorteningAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to dial shortening service:", err)
	}
	defer shorteningConn.Close()
	log.Println("Dialed shortening service successfully!")

	log.Println("Dialing redirection service...")
	redirectionConn, err := grpc.NewClient(config.RedirectionAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to dial redirection service:", err)
	}
	defer redirectionConn.Close()
	log.Println("Dialed redirection service successfully!")

	s := service.NewGatewayService(
		shorteningpb.NewShorteningServiceClient(shorteningConn),
		redirectionpb.NewRedirectionServiceClient(redirectionConn),
	)

	h := handler.NewHandler(*s)

	r := gin.Default()

	handlergen.RegisterHandlers(r, h)

	srv := &http.Server{
		Addr:    net.JoinHostPort("", config.ServerPort),
		Handler: r,
	}

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server is running on port:", config.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()
	<-sigQuit
	log.Println("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server shutdown is successful!")

}

type Config struct {
	ServerPort      string `envconfig:"PORT" required:"true"`
	RedirectionAddr string `envconfig:"REDIRECTION_ADDR" required:"true"`
	ShorteningAddr  string `envconfig:"SHORTENING_ADDR" required:"true"`
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
