package ships

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
)

type externalApi interface {
	GetShip(url string) (*http.Response, error)
}

type apiFunc func()

func (f apiFunc) GetShip(url string) (*http.Response, error) {
	return http.Get(url)
}

type ShipApiRepository struct {
	externalApi
}

func NewShipApiRepository() *ShipApiRepository {
	var apiFunc apiFunc
	return &ShipApiRepository{apiFunc}
}

const API_ENDPOINT = "https://swapi.dev/api/starships/"

type responseShip struct {
	Name          string
	Model         string
	CostInCredits string `json:"cost_in_credits"`
}

func generateShipUrl(shipId string) string {
	return API_ENDPOINT + shipId
}

func (repo ShipApiRepository) GetShip(shipId string) (*domain.Ship, error) {
	response, err := repo.externalApi.GetShip(generateShipUrl(shipId))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ERROR reading ship json from the api due to: %s", err)
	}

	var apiObject responseShip

	err = json.Unmarshal(body, &apiObject)

	if err != nil {
		return nil, fmt.Errorf("ERROR reading ships json from the api due to: %s", err)
	}

	if apiObject.CostInCredits == "unknown" {
		return nil, fmt.Errorf("ERROR this ship does not have a known cost")
	}

	ship, err := domain.NewShipWithId(shipId, apiObject.Name, apiObject.Model, apiObject.CostInCredits)
	if err != nil {
		return nil, fmt.Errorf("ERROR creating ship object due to: %s", err)
	}

	return ship, nil
}
