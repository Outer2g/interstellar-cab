package reservations

import (
	"time"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
	"github.com/google/uuid"
)

type database struct {
	shipOccupation map[int64][]rent
}

type rent struct {
	id        string
	dateFrom  time.Time
	dateTo    time.Time
	userEmail string
	shipId    int64
}

func NewReservationInMemoryDatabase() *database {
	return &database{map[int64][]rent{}}
}

func (repo database) AddReservation(email string, shipId int64, dateFrom, dateTo time.Time) error {
	data := rent{uuid.New().String(), dateFrom, dateTo, email, shipId}

	repo.shipOccupation[shipId] = append(repo.shipOccupation[shipId], data)

	return nil
}

func (repo database) ListReservations(email string) ([]domain.Reservation, error) {
	return nil, nil
}

func (repo database) ShipAvailable(shipId int64, dateFrom, dateTo time.Time) (bool, error) {
	reservationList, isPresent := repo.shipOccupation[shipId]
	if !isPresent {
		// ship has no reservations yet
		return true, nil
	}

	for _, reservation := range reservationList {
		if available := checkReservationDates(reservation.dateFrom, reservation.dateTo, dateFrom, dateTo); !available {
			return false, nil
		}
	}

	return true, nil
}

func checkReservationDates(reservationFrom, reservationTo, newFrom, newTo time.Time) bool {
	if newTo.Before(reservationFrom) || newFrom.After(reservationTo) {
		return true
	}

	return false
}
