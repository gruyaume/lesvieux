package main

import (
	"flag"
	"log"
	"os"

	"github.com/gruyaume/lesvieux/internal/config"
	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/gruyaume/lesvieux/internal/server"
)

func main() {
	log.SetOutput(os.Stderr)
	configFilePtr := flag.String("config", "", "The config file to be provided to the server")
	flag.Parse()
	if *configFilePtr == "" {
		log.Fatalf("Providing a config file is required.")
	}
	conf, err := config.Validate(*configFilePtr)
	if err != nil {
		log.Fatalf("Couldn't validate config file: %s", err)
	}
	dbQueries, err := db.Initialize(conf.DBPath)
	if err != nil {
		log.Fatalf("Couldn't initialize database: %s", err)
	}
	srv, err := server.New(conf.Port, conf.TLS.Cert, conf.TLS.Key, dbQueries)
	if err != nil {
		log.Fatalf("Couldn't create server: %s", err)
	}
	log.Printf("Starting server at %s", srv.Addr)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Server ran into error: %s", err)
	}
}
