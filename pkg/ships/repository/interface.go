package ships

import "github.com/Outer2g/interstellar-cab/pkg/ships"

type ShipsRepository interface {
	GetShip(shipId string) (*ships.Ship, error)
}
