package ships

import (
	"fmt"
	"regexp"
	"strconv"
)

type Ship struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
	Cost  int64  `json:"cost"`
}

func NewShip(url, name, model, cost string) (*Ship, error) {
	// TODO(alex) should not be compiled each time
	r := regexp.MustCompile(`.*/([0-9]+)/`)
	idString := r.FindStringSubmatch(url)

	if len(idString) != 2 {
		return nil, fmt.Errorf("ERROR Could not parse id for ship with url %s", url)
	}

	// The regex already matches for an integer, so no need to re-check for error
	id, _ := strconv.ParseInt(idString[1], 10, 64)

	intCost, err := strconv.ParseInt(cost, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ERROR Could not parse cost for ship with id %d", id)
	}

	return &Ship{id, name, model, intCost}, nil
}

func NewShipWithId(idString, name, model, cost string) (*Ship, error) {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ERROR Could not parse id %s for ship. It is not an integer", idString)
	}

	intCost, err := strconv.ParseInt(cost, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ERROR Could not parse cost %s for ship with id %d", cost, id)
	}

	return &Ship{id, name, model, intCost}, nil
}
