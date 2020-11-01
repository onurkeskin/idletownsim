package sessions

import (
	"github.com/app/server/domain"
)

type IRevokedTokenRepositoryFactory interface {
	New(db domain.IDatabase) IRevokedTokenRepository
}
type IRevokedTokenRepository interface {
	CreateRevokedToken(token IRevokedToken) error
	DeleteExpiredTokens() error
	IsTokenRevoked(id string) bool
}
