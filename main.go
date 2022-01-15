package main

import (
	"context"
	"github.com/SemmiDev/go-song/common/config"
	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/SemmiDev/go-song/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func main() {
	// load the configuration
	config, err := config.Load(".")
	fatal(err, "cannot load config")

	// setup the database with connection pooling
	dbPool, err := pgxpool.Connect(context.Background(), config.DBSource)
	fatal(err, "unable to connect to database")
	defer dbPool.Close()

	// setup the data store
	datastore := db.NewStore(dbPool)

	// setup the server
	s, err := server.New(config, datastore)
	fatal(err, "cannot create server")

	// start the server
	err = s.Start(config.ServerAddress)
	fatal(err, "cannot start server")
}

func fatal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
