package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/readinglist/repository"
	"github.com/wignn/micro-3/readinglist/server"
	"github.com/wignn/micro-3/readinglist/service"
)

type Config struct {
	DSN  string `envconfig:"DATABASE_URL"`
	PORT int    `envconfig:"PORT" default:"50051"`
}

func main() {
	var cfg Config
	fmt.Println("Starting ReadingList Service...")
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}

	var r *repository.PostgresRepository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = repository.NewPostgresRepository(cfg.DSN)
		if err != nil {
			log.Println("Failed to connect to database, retrying...")
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("listening on port", cfg.PORT)
	s := service.NewReadingListService(r)
	log.Fatal(server.ListenGRPC(s, cfg.PORT))
}
