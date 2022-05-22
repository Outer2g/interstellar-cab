package main

import (
	"log"
	"net/http"

	"github.com/Outer2g/interstellar-cab/pkg/service"
	"github.com/gorilla/mux"
)

func simpleResponse(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	listHandler := service.NewListShipsService()

	router.HandleFunc("/register", simpleResponse)
	router.HandleFunc("/listShips", listHandler.HandleListShips)

	router.HandleFunc("/createReservation", simpleResponse)
	router.HandleFunc("/reservations", simpleResponse)

	http.ListenAndServe(":3000", router)
}

func main() {
	log.Printf("Bootstrapping server...")
	handleRequest()
}
