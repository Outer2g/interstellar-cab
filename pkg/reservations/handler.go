package reservations

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	reservations "github.com/Outer2g/interstellar-cab/pkg/reservations/repository"
	ships "github.com/Outer2g/interstellar-cab/pkg/ships/repository"
)

const TO_DAY = 24
const MAX_FREE_USER_COST = 250000

type ReservationHandler struct {
	shipRepository        ships.ShipsRepository
	reservationRepository reservations.ReservationRepository
}

type requestShip struct {
	Id       string    `json:"id"`
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

func NewReservationHandler() *ReservationHandler {
	return &ReservationHandler{ships.NewShipApiRepository(), reservations.NewReservationInMemoryDatabase()}
}

func readRequest(reader io.ReadCloser) (*requestShip, error) {
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("ERROR reading request due to: %s", err)
	}

	var ship requestShip

	err = json.Unmarshal(body, &ship)

	if err != nil {
		return nil, fmt.Errorf("ERROR reading request due to: %s", err)
	}

	return &ship, nil
}

func validateDate(datefrom, dateTo time.Time) error {
	days := int(dateTo.Sub(datefrom).Hours() / TO_DAY)
	log.Println("got days", days)
	if days < 1 || days > 15 {
		return fmt.Errorf("Date not correct, you can only reserve between 1 and 15 days")
	}

	return nil
}

func (u ReservationHandler) HandleNewReservation(rw http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("Email")
	vip, err := strconv.ParseBool(r.Header.Get("Vip"))
	log.Println("new reservation request for ", email, vip)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	requestShip, err := readRequest(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = validateDate(requestShip.DateFrom, requestShip.DateTo)
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		log.Println(err)
		return
	}

	ship, err := u.shipRepository.GetShip(requestShip.Id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Ship not found due to", err)
		return
	}

	if !vip && ship.Cost > MAX_FREE_USER_COST {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("User cannot make a reservation for that ship. He is not premium")
		return
	}

	isAvailable, err := u.reservationRepository.ShipAvailable(ship.Id, requestShip.DateFrom, requestShip.DateTo)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not check ship availability")
		return
	}

	if !isAvailable {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Ship not available for those dates")
		return
	}

	err = u.reservationRepository.AddReservation(email, ship.Id, requestShip.DateFrom, requestShip.DateTo)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not make ship reservation")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (handler ReservationHandler) HandleListReservations(rw http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("Email")
	vip, _ := strconv.ParseBool(r.Header.Get("Vip"))
	log.Println("new list reservation request for", email, vip)

	list := handler.reservationRepository.ListReservations(email)
	log.Println("User has reservations:", list)

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(list)
}
