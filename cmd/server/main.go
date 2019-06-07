package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/volmedo/pAPI/pkg/restapi"
	"github.com/volmedo/pAPI/pkg/service"
)

func main() {
	var dbHost, dbUser, dbPass, dbName, migrationsPath string
	var port, dbPort int
	var rps int64
	flag.IntVar(&port, "port", 8080, "Port where the server is listening for connections.")
	flag.Int64Var(&rps, "rps", 100, "Rate limit expressed in requests per second (per client)")

	flag.StringVar(&dbHost, "dbhost", "localhost", "Address of the server that hosts the DB")
	flag.IntVar(&dbPort, "dbport", 5432, "Port where the DB server is listening for connections")
	flag.StringVar(&dbUser, "dbuser", "postgres", "User to use when accessing the DB")
	flag.StringVar(&dbPass, "dbpass", "postgres", "Password to use when accessing the DB")
	flag.StringVar(&dbName, "dbname", "postgres", "Name of the DB to connect to")
	flag.StringVar(&migrationsPath, "migrations", "./migrations", "Path to the folder that contains the migration files")

	flag.Parse()

	// Setup DB
	dbConf := &service.DBConfig{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Pass:           dbPass,
		Name:           dbName,
		MigrationsPath: migrationsPath,
	}
	db, err := service.NewDB(dbConf)
	if err != nil {
		log.Panicf("Unable to configure DB connection: %v", err)
	}

	testRepo, err := service.NewDBPaymentRepository(db, dbName, migrationsPath)
	if err != nil {
		log.Panicf("Unable to create DB repo: %v", err)
	}

	ps := &service.PaymentsService{Repo: testRepo}

	apiHandler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: ps,
		Logger:      log.Printf,
	})
	if err != nil {
		log.Panicf("Error creating main API handler: %v", err)
	}

	apiHandler, prometheusHandler := newMeasuredHandler(apiHandler)
	apiHandler, err = newRateLimitedHandler(rps, apiHandler)
	if err != nil {
		log.Panicf("Error creating rate limiter middleware: %v", err)
	}
	apiHandler = newRecoverableHandler(apiHandler)

	mux := http.NewServeMux()
	mux.Handle("/metrics", prometheusHandler)
	mux.Handle("/", apiHandler)

	log.Printf("Starting server, accepting requests on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Panicf("Error while serving: %v", err)
	}
}
