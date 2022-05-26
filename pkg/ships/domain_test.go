package ships

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShip(t *testing.T) {
	t.Run("Should create proper ship object when all data is correct", func(t *testing.T) {
		expectedShip := &Ship{1, "test", "testModel", 2000}
		ship, err := NewShip(aUrlWithId(1), "test", "testModel", "2000")

		assert.EqualValues(t, nil, err)
		assert.EqualValues(t, expectedShip, ship)
	})

	t.Run("Should return error when invalid id in url", func(t *testing.T) {
		_, err := NewShip("arandomurl231", "test", "testModel", "2000")

		assert.Containsf(t, err.Error(), "ERROR Could not parse id for ship with url", err.Error())
	})

	t.Run("Should return error when invalid cost", func(t *testing.T) {
		_, err := NewShip(aUrlWithId(1), "test", "testModel", "notACost")

		assert.Containsf(t, err.Error(), "ERROR Could not parse cost for ship with id", err.Error())
	})
}

func aUrlWithId(id int64) string {
	return fmt.Sprintf("https://example.org/api/starships/%d/", id)
}
