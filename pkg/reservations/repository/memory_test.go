package reservations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestReservationsDatabase() *database {
	from := time.Date(2022, 05, 22, 0, 0, 0, 0, time.UTC)
	to := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)
	creation  := time.Date(2022, 05, 20, 0, 0, 0, 0, time.UTC)

	reservation := Reservation{"21598182", from, to, "existing@email.com", 12, creation}

	shipOccupation := map[int64][]Reservation{12: {reservation}}
	userOccupation := map[string][]Reservation{"existing@email.com": {reservation}}
	return &database{shipOccupation, userOccupation}
}

func TestAddReservation(t *testing.T) {
	repository := newTestReservationsDatabase()

	t.Run("Should add new reservation", func(t *testing.T) {
		from := time.Date(2022, 05, 22, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)
		email := "anonexistent@email.com"

		err := repository.AddReservation(email, 16, from, to)

		assert.Nil(t, err)
		assert.EqualValues(t, 1, len(repository.shipOccupation[16]))
		assert.EqualValues(t, 1, len(repository.userOccupation[email]))
		shipReservation := repository.shipOccupation[16][0]
		userReservation := repository.userOccupation[email][0]

		assert.EqualValues(t, shipReservation, userReservation)
		assert.EqualValues(t, shipReservation.DateFrom, from)
		assert.EqualValues(t, shipReservation.DateTo, to)
		assert.EqualValues(t, shipReservation.ShipId, 16)
		assert.EqualValues(t, shipReservation.UserEmail, "anonexistent@email.com")
	})
}

func TestListUserReservation(t *testing.T) {

	t.Run("Should list existing reservations", func(t *testing.T) {
		repository := newTestReservationsDatabase()
		from := time.Date(2022, 05, 22, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)
		creation := time.Date(2022, 05, 20, 0, 0, 0, 0, time.UTC)
		email := "existing@email.com"

		result := repository.ListReservations(email)

		assert.EqualValues(t, []Reservation{{"21598182", from, to, "existing@email.com", 12, creation}}, result)
	})

	t.Run("Should return empty list when no reservations", func(t *testing.T) {
		repository := newTestReservationsDatabase()
		email := "anonexisting@email.com"

		result := repository.ListReservations(email)

		assert.EqualValues(t, []Reservation{}, result)
	})
}

func TestShipAvailable(t *testing.T) {
	repository := newTestReservationsDatabase()

	t.Run("Should return true when the ship has no reservations", func(t *testing.T) {
		from := time.Date(2022, 05, 22, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(16, from, to)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("Should return true when the ship is available for those dates when dates are BEFORE the reserved ones", func(t *testing.T) {
		// Ship 12 is reserved for days 22-25
		from := time.Date(2022, 05, 19, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 21, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(12, from, to)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("Should return true when the ship is available for those dates when dates are AFTER the reserved ones", func(t *testing.T) {
		// Ship 12 is reserved for days 22-25
		from := time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 30, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(12, from, to)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("Should return false when the ship is not available for those dates when it overlaps in the low range", func(t *testing.T) {
		// Ship 12 is reserved for days 22-25
		from := time.Date(2022, 05, 20, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 22, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(12, from, to)

		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("Should return false when the ship is not available for those dates when it overlaps in the upper range", func(t *testing.T) {
		// Ship 12 is reserved for days 22-25
		from := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(12, from, to)

		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("Should return false when the ship is not available for those dates when it overlaps in the upper range", func(t *testing.T) {
		// Ship 12 is reserved for days 22-25
		from := time.Date(2022, 05, 25, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC)

		result, err := repository.ShipAvailable(12, from, to)

		assert.Nil(t, err)
		assert.False(t, result)
	})
}
