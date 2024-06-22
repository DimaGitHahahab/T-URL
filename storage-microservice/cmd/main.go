package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"storage/internal/repository/postgres"
	"storage/internal/repository/redis"
	"storage/internal/server"
	"storage/internal/service"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	redislib "github.com/redis/go-redis/v9"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Loading config...")
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Loaded config successfully!")

	ctx := context.Background()

	log.Println("Setting up pgx pool...")
	pool, err := postgres.SetupPgxPool(ctx, config.PostgresConnString)
	if err != nil {
		log.Fatal("Failed to setup pgx pool:", err)
	}
	log.Println("Pgx pool setup successfully!")

	log.Println("Processing postgres migration...")
	err = postgres.ProcessMigration(config.MigrationPath, config.PostgresConnString)
	if err != nil {
		log.Fatal("Failed to process postgres migration:", err)
	}
	log.Println("Postgres migration processed successfully!")

	repo := postgres.New(pool)

	redisClient := redislib.NewClient(&redislib.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	c := redis.New(redisClient)
	s := service.New(repo, c)

	srv := server.New(s)

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server is running on port:", config.Port)
		log.Fatal(srv.Run(config.Port))
	}()

	<-sigQuit
	log.Println("Gracefully shutting down server...")
	srv.Stop()

	log.Println("Server shutdown is successful!")
}

type Config struct {
	Port string `envconfig:"PORT" required:"true"`

	PostgresConnString string `envconfig:"POSTGRES_URL" required:"true"`
	MigrationPath      string `envconfig:"MIGRATION_PATH" required:"true"`

	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"true"`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
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
