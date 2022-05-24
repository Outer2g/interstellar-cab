package ships

import (
	"github.com/Outer2g/interstellar-cab/pkg/domain"
)

type ShipsRepository interface {
	GetShip(shipId string) (*domain.Ship, error)
}
