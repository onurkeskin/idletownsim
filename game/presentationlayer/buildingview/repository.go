package buildingview

import (
	"errors"
	"fmt"
	buildingviewdomain "github.com/app/game/presentationlayer/buildingview/domain"
	"github.com/app/server/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const BuildingViewCollection string = "buildingview"

type BuildingViewRepository struct {
	DB domain.IDatabase
}

// CreateBuildingView Insert new BuildingView document into the database
func (repo *BuildingViewRepository) CreateBuildingView(_b buildingviewdomain.IBuildingView) error {
	b := _b.(*BuildingView)
	if b.ID.Hex() == "" {
		b.ID = bson.NewObjectId()
	}
	return repo.DB.Insert(BuildingViewCollection, b)
}

// GetUsers Get list of users
func (repo *BuildingViewRepository) GetBuildingsView() buildingviewdomain.IBuildingView {
	b := BuildingsView{}
	err := repo.DB.FindAll(BuildingViewCollection, nil, &b, 50, "")
	if err != nil {
		return nil
	}
	return &b
}

// GetUser Get user specified by the id
func (repo *BuildingViewRepository) GetBuildingViewById(id string) (buildingviewdomain.IBuildingView, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var b BuildingView
	err := repo.DB.FindOne(BuildingViewCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &b)
	return &b, err
}

func (repo *BuildingViewRepository) GetBuildingViewsByBuildingViewId(bid string) (buildingviewdomain.IBuildingsView, error) {
	if !bson.IsObjectIdHex(bid) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", bid))
	}

	err := repo.DB.EnsureIndex(BuildingViewCollection, mgo.Index{
		Key: []string{
			"$text:globalidentifier",
		},
		Background: true,
		Sparse:     false,
	})
	q := domain.Query{"globalidentifier": bson.ObjectIdHex(bid)}
	var bs BuildingsView
	err = repo.DB.FindAll(BuildingViewCollection, q, &bs, 999999, "")

	return bs, err

}
