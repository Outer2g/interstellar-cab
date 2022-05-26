package reservations

import (
	"time"
)

type ReservationRepository interface {
	AddReservation(email string, shipId int64, dateFrom, dateTo time.Time) error
	ListReservations(email string) []Reservation
	ShipAvailable(shipId int64, dateFrom, dateTo time.Time) (bool, error)
}
