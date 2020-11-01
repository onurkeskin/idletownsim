package resourceview

import (
	"errors"
	"fmt"
	ResourceViewdomain "github.com/app/game/presentationlayer/resourceview/domain"
	"github.com/app/server/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ResourceViewCollection string = "resourceview"

type ResourceViewRepository struct {
	DB domain.IDatabase
}

// CreateResourceView Insert new ResourceView document into the database
func (repo *ResourceViewRepository) CreateResourceView(_r ResourceViewdomain.IResourceView) error {
	r := _r.(*ResourceView)
	if r.ID.Hex() == "" {
		r.ID = bson.NewObjectId()
	}
	return repo.DB.Insert(ResourceViewCollection, r)
}

// GetUsers Get list of users
func (repo *ResourceViewRepository) GetResourceViews() ResourceViewdomain.IResourceView {
	resources := ResourcesView{}
	err := repo.DB.FindAll(ResourceViewCollection, nil, &resources, 50, "")
	if err != nil {
		return nil
	}
	return &resources
}

// GetUser Get user specified by the id
func (repo *ResourceViewRepository) GetResourceViewById(id string) (ResourceViewdomain.IResourceView, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var res ResourceView
	err := repo.DB.FindOne(ResourceViewCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &res)
	return &res, err
}

func (repo *ResourceViewRepository) GetResourceViewsByResourceViewId(bid string) (ResourceViewdomain.IResourcesView, error) {
	if !bson.IsObjectIdHex(bid) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", bid))
	}

	err := repo.DB.EnsureIndex(ResourceViewCollection, mgo.Index{
		Key: []string{
			"$text:globalidentifier",
		},
		Background: true,
		Sparse:     false,
	})
	q := domain.Query{"globalidentifier": bson.ObjectIdHex(bid)}
	var resources ResourcesView
	err = repo.DB.FindAll(ResourceViewCollection, q, &resources, 999999, "")

	return &resources, err

}
