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
	var port, rps, dbPort int
	flag.IntVar(&port, "port", 8080, "Port where the server is listening for connections.")
	flag.IntVar(&rps, "rps", 100, "Rate limit expressed in requests per second (per client)")

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
		panic(fmt.Sprintf("Unable to configure DB connection: %v", err))
	}

	testRepo, err := service.NewDBPaymentRepository(db, dbName, migrationsPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to create DB repo: %v", err))
	}

	ps := &service.PaymentsService{Repo: testRepo}

	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: ps,
		Logger:      log.Printf,
	})
	if err != nil {
		panic(fmt.Sprintf("Error creating main API handler: %v", err))
	}

	rateLimitedHandler, err := newRateLimitedHandler(int64(rps), handler)
	if err != nil {
		panic(fmt.Sprintf("Error creating rate limiter middleware: %v", err))
	}

	log.Printf("Starting server, accepting requests on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), rateLimitedHandler))
}
