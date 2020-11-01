package sessions

import (
	. "github.com/app/sessions/domain"

	"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User collection name
const RevokedTokenCollections string = "revoked_tokens"

func NewRevokedTokenRepositoryFactory() IRevokedTokenRepositoryFactory {
	return &RevokedTokenRepositoryFactory{}
}

type RevokedTokenRepositoryFactory struct{}

func (factory *RevokedTokenRepositoryFactory) New(db domain.IDatabase) IRevokedTokenRepository {
	return &RevokedTokenRepository{db}
}

type RevokedTokenRepository struct {
	DB domain.IDatabase
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) CreateRevokedToken(token IRevokedToken) error {
	t := token.(*RevokedToken)
	t.RevokedDate = time.Now()
	return repo.DB.Insert(RevokedTokenCollections, t)
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) DeleteExpiredTokens() error {
	return repo.DB.RemoveAll(RevokedTokenCollections, domain.Query{
		"exp": domain.Query{
			"$lt": time.Now(),
		},
	})
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) IsTokenRevoked(id string) bool {
	if !bson.IsObjectIdHex(id) {
		return false
	}
	return repo.DB.Exists(RevokedTokenCollections, domain.Query{
		"_id": bson.ObjectIdHex(id),
	})
}
