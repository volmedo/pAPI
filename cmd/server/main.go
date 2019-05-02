package main

import (
	"log"
	"net/http"

	"github.com/volmedo/pAPI/pkg/impl"
	"github.com/volmedo/pAPI/pkg/restapi"
)

const (
	serverURL  = "http://localhost"
	serverPort = "8080"
	apiRoot    = serverURL + ":" + serverPort + "/v1/"
)

func main() {
	p := &impl.PaymentsAPI{}

	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: p,
		Logger:      log.Printf,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server, API can be consumed at %s\n", apiRoot)

	log.Fatal(http.ListenAndServe(":"+serverPort, handler))
}
