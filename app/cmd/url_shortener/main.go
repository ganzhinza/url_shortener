package main

//TODO: units

import (
	"context"
	"log"
	"net/http"
	"os"
	"url_shortener/pkg/api"
	"url_shortener/pkg/config"
	"url_shortener/pkg/db/memdb"
	"url_shortener/pkg/db/pgsql"

	"github.com/jackc/pgx/v5/pgxpool"
)

const URLSIZE = 10

type server struct {
	httpServer *http.Server
	api        *api.API
}

func main() {
	srv := new(server)

	cfg := config.MustLoad()
	switch cfg.StorageType {
	case "memdb":
		srv.api = api.New(cfg.HostName, memdb.NewDB(URLSIZE))
	case "pgsql":
		pool := InitializePostgres()
		defer pool.Close()
		srv.api = api.New(cfg.HostName, pgsql.New(pool, URLSIZE))
	}

	srv.httpServer = &http.Server{
		Addr:         cfg.Address,
		Handler:      srv.api.Router(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func InitializePostgres() *pgxpool.Pool {
	ctx := context.Background()
	usr := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	db, err := pgxpool.New(ctx, "postgres://"+usr+":"+pwd+"@db:5432/url_shortener?sslmode=disable")

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbCreation, err := os.ReadFile("../../config/schema.sql")
	if err != nil {
		log.Fatalf("Unable to read schema: %v\n", err)
	}
	_, err = db.Exec(ctx, string(dbCreation))
	if err != nil {
		log.Fatalf("Unable to create tables: %v\n", err)
	}
	return db
}
