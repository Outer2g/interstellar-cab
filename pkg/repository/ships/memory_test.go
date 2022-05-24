package ships

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
	"github.com/stretchr/testify/assert"
)

type testApi struct {
	t            *testing.T
	expectedUrl  string
	mockResponse string
	mockError    error
}

func (f testApi) GetShip(url string) (*http.Response, error) {
	assert.EqualValues(f.t, f.expectedUrl, url)
	if f.mockError != nil {
		return nil, f.mockError
	}

	body := ioutil.NopCloser(bytes.NewBufferString(f.mockResponse))

	return &http.Response{Body: body}, nil
}

func newTestListShipService(t *testing.T, expectedUrl, shipResponse string, err error) *ShipApiRepository {
	return &ShipApiRepository{testApi{t, expectedUrl, shipResponse, err}}
}
func TestGetShip(t *testing.T) {
	t.Run("Should return ship information", func(t *testing.T) {
		service := newTestListShipService(t, aShipUrl("2"), aShipResponseJson(), nil)
		result, err := service.GetShip("2")

		assert.Nil(t, err)
		assert.Equal(t, aShip(), result)
	})

	t.Run("Should return error when api error", func(t *testing.T) {
		service := newTestListShipService(t, aShipUrl("2"), "", fmt.Errorf("Test error"))
		result, err := service.GetShip("2")

		assert.NotNil(t, err)
		assert.Containsf(t, err.Error(), "Test error", err.Error())
		assert.Nil(t, result)
	})

	// TODO more validation tests
}

func aShip() *domain.Ship {
	return &domain.Ship{2, "CR90 corvette", "CR90 corvette", 3500000}
}

func aShipUrl(id string) string {
	return "https://swapi.dev/api/starships/" + id
}

func aShipResponseJson() string {
	return `{
		"name": "CR90 corvette", 
		"model": "CR90 corvette", 
		"manufacturer": "Corellian Engineering Corporation", 
		"cost_in_credits": "3500000", 
		"length": "150", 
		"max_atmosphering_speed": "950", 
		"crew": "30-165", 
		"passengers": "600", 
		"cargo_capacity": "3000000", 
		"consumables": "1 year", 
		"hyperdrive_rating": "2.0", 
		"MGLT": "60", 
		"starship_class": "corvette", 
		"pilots": [], 
		"films": [
			"https://swapi.dev/api/films/1/", 
			"https://swapi.dev/api/films/3/", 
			"https://swapi.dev/api/films/6/"
		], 
		"created": "2014-12-10T14:20:33.369000Z", 
		"edited": "2014-12-20T21:23:49.867000Z", 
		"url": "https://swapi.dev/api/starships/2/"
	}`
}
