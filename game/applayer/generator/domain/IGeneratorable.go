package domain

import (
	"github.com/app/game/applayer/general/domain"
)

type IGeneratorable interface {
	GetTargetResources() []string
	GetScheme() domain.IValScheme
	GetPriority() int64
}
