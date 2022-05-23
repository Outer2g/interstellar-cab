package reservations

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
	"github.com/stretchr/testify/assert"
)

type shipRepository struct {
	mockShip *domain.Ship
}

func (repo shipRepository) GetShip(shipId string) *domain.Ship {
	return repo.mockShip
}

type reservationTestRepository struct {
	mockAvailability      bool
	mockAvailabilityError error
	mockReservationError  error
}

func (repo reservationTestRepository) AddReservation(email string, shipId int64, dateFrom, dateTo time.Time) error {
	return repo.mockReservationError
}
func (repo reservationTestRepository) ListReservations(email string) ([]domain.Reservation, error) {
	//not used
	return nil, nil
}
func (repo reservationTestRepository) ShipAvailable(shipId int64, dateFrom, dateTo time.Time) (bool, error) {
	return repo.mockAvailability, repo.mockAvailabilityError
}

func newTestReservationHandler(mockShip *domain.Ship, mockAvailability bool, mockAvailabilityError, mockReservationError error) *ReservationHandler {
	ships := shipRepository{mockShip}
	reservations := reservationTestRepository{mockAvailability, mockAvailabilityError, mockReservationError}
	return &ReservationHandler{ships, reservations}
}

func TestHandleShipReservations(t *testing.T) {
	t.Run("Should make reservation", func(t *testing.T) {
		handler := newTestReservationHandler(aShipWithId(4), true, nil, nil)
		req := httptest.NewRequest("POST", "/reservation", aShipInJson())
		req.Header.Set("Email", "existing@email.com")
		req.Header.Set("Vip", "true")
		recorder := httptest.NewRecorder()

		handler.HandleNewReservation(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 200, result.StatusCode)
	})

	t.Run("Should return error when trying to make a reservation for over 15 days ", func(t *testing.T) {
		handler := newTestReservationHandler(aShipWithId(4), true, nil, nil)
		req := httptest.NewRequest("POST", "/reservation", aShipWithMoreThan15DaysAsDatesJson())
		req.Header.Set("Email", "existing@email.com")
		req.Header.Set("Vip", "true")
		recorder := httptest.NewRecorder()

		handler.HandleNewReservation(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})

	t.Run("Should return error when trying to make a reservation  with invalid dates ", func(t *testing.T) {
		handler := newTestReservationHandler(aShipWithId(4), true, nil, nil)
		req := httptest.NewRequest("POST", "/reservation", aShipWithMoreThan15DaysAsDatesJson())
		req.Header.Set("Email", "existing@email.com")
		req.Header.Set("Vip", "true")
		recorder := httptest.NewRecorder()

		handler.HandleNewReservation(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})

	t.Run("Should return error when user not premium and cost of ship over threshold", func(t *testing.T) {
		handler := newTestReservationHandler(aShipWithId(10), true, nil, nil)
		req := httptest.NewRequest("POST", "/reservation", aShipWithMoreThan15DaysAsDatesJson())
		req.Header.Set("Email", "existing@email.com")
		req.Header.Set("Vip", "false")
		recorder := httptest.NewRecorder()

		handler.HandleNewReservation(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})

	//TODO add tests that are outside happy path
}
func aShipWithWrongDates() *bytes.Buffer {
	var jsonData = []byte(`{"id": "14","date_from": "2022-07-22T16:00:00.000Z","date_to": "2022-05-22T15:00:00.000Z"}`)
	return bytes.NewBuffer(jsonData)
}

func aShipWithMoreThan15DaysAsDatesJson() *bytes.Buffer {
	var jsonData = []byte(`{
		"id": "14",
		"date_from": "2022-05-22T15:00:00.000Z",
		"date_to": "2022-07-22T16:00:00.000Z"
	}`)
	return bytes.NewBuffer(jsonData)
}

func aShipInJson() *bytes.Buffer {
	var jsonData = []byte(`{"id": "12","date_from": "2022-05-22T15:00:00.000Z","date_to": "2022-05-24T16:00:00.000Z"}`)
	return bytes.NewBuffer(jsonData)
}

func aShipWithId(id int64) *domain.Ship {
	return &domain.Ship{id, "test", "testModel", 4000000000}
}
