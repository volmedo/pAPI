package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/namsral/flag"

	"github.com/volmedo/pAPI/pkg/restapi"
	"github.com/volmedo/pAPI/pkg/service"
)

func main() {
	var dbHost, dbUser, dbPass, dbName, migrationsPath string
	var port, dbPort int
	var rps int64

	// Use "PAPI" as prefix for env variables to avoid potential clashes
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "PAPI", flag.ExitOnError)

	fs.IntVar(&port, "port", 8080, "Port where the server is listening for connections.")
	fs.Int64Var(&rps, "rps", 100, "Rate limit expressed in requests per second (per client)")

	fs.StringVar(&dbHost, "dbhost", "localhost", "Address of the server that hosts the DB")
	fs.IntVar(&dbPort, "dbport", 5432, "Port where the DB server is listening for connections")
	fs.StringVar(&dbUser, "dbuser", "postgres", "User to use when accessing the DB")
	fs.StringVar(&dbPass, "dbpass", "postgres", "Password to use when accessing the DB")
	fs.StringVar(&dbName, "dbname", "postgres", "Name of the DB to connect to")
	fs.StringVar(&migrationsPath, "migrations", "./migrations", "Path to the folder that contains the migration files")

	// Ignore errors; fs is set for ExitOnError
	_ = fs.Parse(os.Args[1:])

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)

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
		logger.Panicf("Unable to configure DB connection: %v", err)
	}

	testRepo, err := service.NewDBPaymentRepository(db, dbName, migrationsPath)
	if err != nil {
		logger.Panicf("Unable to create DB repo: %v", err)
	}

	ps := &service.PaymentsService{
		Repo:   testRepo,
		Logger: logger,
	}

	apiHandler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: ps,
		Logger:      logger.Printf,
	})
	if err != nil {
		logger.Panicf("Error creating main API handler: %v", err)
	}

	apiHandler, prometheusHandler := newMeasuredHandler(apiHandler)
	apiHandler, err = newRateLimitedHandler(rps, apiHandler)
	if err != nil {
		logger.Panicf("Error creating rate limiter middleware: %v", err)
	}
	apiHandler = newRecoverableHandler(apiHandler)

	mux := http.NewServeMux()
	mux.Handle("/health", newHealthHandler(db))
	mux.Handle("/metrics", prometheusHandler)
	mux.Handle("/", apiHandler)

	logger.Printf("Starting server, accepting requests on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		logger.Panicf("Error while serving: %v", err)
	}
}
