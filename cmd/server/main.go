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
	port := flag.Int("port", 8080, "Port where the server is listening for connections.")
	flag.Parse()

	ps := &service.PaymentsService{
		Repo: service.NewPaymentRepository(),
	}

	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: ps,
		Logger:      log.Printf,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server, accepting requests on port %d\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
