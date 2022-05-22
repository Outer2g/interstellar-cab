package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Outer2g/interstellar-cab/pkg/domain"
)

const API_ENDPOINT = "https://swapi.dev/api/starships"

type shipResponse struct {
	Url   string
	Name  string
	Model string
	Cost  string `json:"cost_in_credits"`
}

type apiResponse struct {
	Next    string         `json:"next"`
	Results []shipResponse `json:"results"`
}

type externalApi interface {
	getShips() (resp *http.Response, err error)
}

type apiFunc func()

func (f apiFunc) getShips() (resp *http.Response, err error) {
	return http.Get(API_ENDPOINT)
}

func NewListShipsService() *ListShipsService {
	var impl apiFunc

	return &ListShipsService{impl}
}

type ListShipsService struct {
	externalApi
}

func (s *ListShipsService) getShipsFromApi() (*apiResponse, error) {
	// TODO(alex) retrieve all ships using the next field
	response, err := s.getShips()

	if err != nil {
		return nil, fmt.Errorf("ERROR reading ships from the api due to: %s", err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var apiObject apiResponse

	err = json.Unmarshal(body, &apiObject)

	if err != nil {
		return nil, fmt.Errorf("ERROR reading ships json from the api due to: %s", err)
	}

	log.Println("Got response", apiObject)

	if apiObject.Results == nil {
		return nil, fmt.Errorf("Got Empty results from the API call: %s", err)
	}

	return &apiObject, nil
}

func (s *ListShipsService) HandleListShips(rw http.ResponseWriter, r *http.Request) {

	apiObject, err := s.getShipsFromApi()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	result := []domain.Ship{}
	for _, responseShip := range apiObject.Results {
		ship, err := domain.NewShip(responseShip.Url, responseShip.Name, responseShip.Model, responseShip.Cost)
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, *ship)
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}

	rw.Write(data)
}
