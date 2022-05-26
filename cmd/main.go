package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Outer2g/interstellar-cab/pkg/auth"
	"github.com/Outer2g/interstellar-cab/pkg/reservations"
	"github.com/Outer2g/interstellar-cab/pkg/ships"
	"github.com/Outer2g/interstellar-cab/pkg/user"
	"github.com/gorilla/mux"
)

func simpleResponse(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	listHandler := ships.NewListShipsHandler()
	userHandler := user.NewUserAuth()
	reservationHandler := reservations.NewReservationHandler()

	router.HandleFunc("/login", userHandler.HandleLoginUser).Methods(http.MethodPost)
	router.HandleFunc("/signup", userHandler.HandleSignupUser).Methods(http.MethodPost)

	router.HandleFunc("/ships", listHandler.HandleListShips)

	router.HandleFunc("/reservations", auth.CheckAuth(reservationHandler.HandleNewReservation)).Methods(http.MethodPost)
	router.HandleFunc("/reservations", auth.CheckAuth(reservationHandler.HandleListReservations))

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
