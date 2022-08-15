package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nobia/simplebank/api"
	db "github.com/nobia/simplebank/db/sqlc"
	"github.com/nobia/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	log.Fatal(server.Start(config.ServerAddress))
}
