package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Outer2g/interstellar-cab/pkg/auth"
	"github.com/Outer2g/interstellar-cab/pkg/reservations"
	"github.com/Outer2g/interstellar-cab/pkg/service"
	"github.com/gorilla/mux"
)

func simpleResponse(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	listHandler := service.NewListShipsService()
	authService := auth.NewUserAuth()
	reser := reservations.NewReservationHandler()

	router.HandleFunc("/login", authService.HandleLoginUser).Methods(http.MethodPost)
	router.HandleFunc("/signup", authService.HandleSignupUser).Methods(http.MethodPost)

	router.HandleFunc("/listShips", listHandler.HandleListShips)

	router.HandleFunc("/createReservation", authService.CheckAuth(reser.HandleNewReservation)).Methods(http.MethodPost)
	router.HandleFunc("/reservations", authService.CheckAuth(simpleResponse))

	http.ListenAndServe(":3000", router)
}

func main() {
	log.Println("Checking environment variables...")

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	log.Printf("Bootstrapping server...")
	handleRequest()
}
