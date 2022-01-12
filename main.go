package main

import (
	"context"
	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/SemmiDev/go-song/server"
	"github.com/SemmiDev/go-song/util"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config, err := util.LoadEnv(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	store := db.NewStore(dbpool)
	s, err := server.New(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = s.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
