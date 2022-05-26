package reservations

import (
	"time"

	"github.com/google/uuid"
)

type database struct {
	shipOccupation map[int64][]Reservation
	userOccupation map[string][]Reservation
}

type Reservation struct {
	Id        string
	DateFrom  time.Time
	DateTo    time.Time
	UserEmail string
	ShipId    int64
}

func NewReservationInMemoryDatabase() *database {
	return &database{map[int64][]Reservation{}, map[string][]Reservation{}}
}

func (repo database) AddReservation(email string, shipId int64, dateFrom, dateTo time.Time) error {
	data := Reservation{uuid.New().String(), dateFrom, dateTo, email, shipId}

	repo.shipOccupation[shipId] = append(repo.shipOccupation[shipId], data)
	repo.userOccupation[email] = append(repo.userOccupation[email], data)

	return nil
}

func (repo database) ListReservations(email string) []Reservation {
	reservations, present := repo.userOccupation[email]

	if !present {
		return []Reservation{}
	}

	return reservations
}

func (repo database) ShipAvailable(shipId int64, dateFrom, dateTo time.Time) (bool, error) {
	reservationList, isPresent := repo.shipOccupation[shipId]
	if !isPresent {
		// ship has no reservations yet
		return true, nil
	}

	for _, reservation := range reservationList {
		if available := checkReservationDates(reservation.DateFrom, reservation.DateTo, dateFrom, dateTo); !available {
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
