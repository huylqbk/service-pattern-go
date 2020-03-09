package interfaces

import (
	"github.com/huylqbk/service-pattern-go/models"
)

type IPlayerRepository interface {
	GetPlayerByName(name string) (models.PlayerModel, error)
}
