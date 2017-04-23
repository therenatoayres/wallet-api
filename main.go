package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/therenatoayres/wallet-api/router"
)

const (
	listeningPort string = "8144"
)

func main() {
	log.Printf("Server started")

	rout := mux.
		NewRouter().
		StrictSlash(true)
	rout = router.Router(rout)
	log.Fatal(http.ListenAndServe(":"+listeningPort, rout))
}
