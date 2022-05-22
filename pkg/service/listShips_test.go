package service

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testApi struct {
	mockResponse string
	mockError    error
}

func (f testApi) getShips() (resp *http.Response, err error) {
	if f.mockError != nil {
		return nil, f.mockError
	}

	body := ioutil.NopCloser(bytes.NewBufferString(f.mockResponse))

	return &http.Response{Body: body}, err
}

func newTestListShipService(value string, err error) *ListShipsService {
	impl := testApi{value, err}
	return &ListShipsService{impl}
}
func TestHandeListShips(t *testing.T) {
	t.Run("Should return list of Ships when no pagination", func(t *testing.T) {
		service := newTestListShipService(anApiResponse("null", aShip(1, "test", "testModel", 100), aShip(2, "test2", "testModel", 200)), nil)
		req := httptest.NewRequest("GET", "/listShips", nil)
		recorder := httptest.NewRecorder()
		service.HandleListShips(recorder, req)
		result := recorder.Result()

		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 200, result.StatusCode)
		assert.EqualValues(t, "[{\"id\":1,\"name\":\"test\",\"model\":\"testModel\",\"cost\":100},{\"id\":2,\"name\":\"test2\",\"model\":\"testModel\",\"cost\":200}]", string(body))
	})

	t.Run("Should return error when error on the API", func(t *testing.T) {
		service := newTestListShipService("", errors.New("Error in the API"))
		req := httptest.NewRequest("GET", "/listShips", nil)
		recorder := httptest.NewRecorder()
		service.HandleListShips(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})

	t.Run("Should return error when invalid json", func(t *testing.T) {
		service := newTestListShipService("notajson", nil)
		req := httptest.NewRequest("GET", "/listShips", nil)
		recorder := httptest.NewRecorder()
		service.HandleListShips(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})

	t.Run("Should return error when json is retrieved without proper Result field", func(t *testing.T) {
		service := newTestListShipService(aBrokenApiResponse(), nil)
		req := httptest.NewRequest("GET", "/listShips", nil)
		recorder := httptest.NewRecorder()
		service.HandleListShips(recorder, req)
		result := recorder.Result()

		assert.EqualValues(t, 503, result.StatusCode)
	})
}

func aBrokenApiResponse() string {
	return "{\"next\":null,\"previous\":null}"
}

func anApiResponse(next string, ships ...string) string {
	return fmt.Sprintf("{\"next\":%s,\"previous\":null,\"results\":[%s]}", next, strings.Join(ships, ","))
}

func aShip(id int64, name, model string, cost int64) string {
	return fmt.Sprintf("{\"url\":\"https://example.org/api/starships/%d/\",\"name\":\"%s\",\"model\":\"%s\",\"cost_in_credits\":\"%d\"}", id, name, model, cost)
}
