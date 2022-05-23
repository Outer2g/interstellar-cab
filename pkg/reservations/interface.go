package reservations

import (
	"time"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
)

type ReservationRepository interface {
	AddReservation(email string, shipId int64, dateFrom, dateTo time.Time) error
	ListReservations(email string) ([]domain.Reservation, error)
	ShipAvailable(shipId int64, dateFrom, dateTo time.Time) (bool, error)
}
