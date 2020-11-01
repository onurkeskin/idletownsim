package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/app/game/applayer/building"
	upgrade "github.com/app/game/applayer/upgrade"
	"github.com/app/game/presentationlayer"
	"github.com/app/users"
	"github.com/datastructures/datastructures"
	"github.com/gorilla/mux"
	"time"
	//building "github.com/app/game/applayer/building"
	"github.com/app/game/applayer/gamemap"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	"github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	maps "github.com/app/map"
	mapmodels "github.com/app/map/mapmodels"
	"gopkg.in/mgo.v2/bson"
	//buildingdomain "github.com/app/game/applayer/building/domain"
	buildinginstance "github.com/app/game/applayer/buildinginstance"
	game "github.com/app/game/applayer/game"
	gamedomain "github.com/app/game/applayer/game/domain"
	//maps "github.com/app/map"

	//"github.com/app/server/domain"
	//"gopkg.in/mgo.v2/bson"
	//"log"
	"net/http"
)

/*
type GetGamesForUserRequest_v0 struct {
	Game    game.GameEnvironment `json:"game"`
	Success bool                 `json:"success"`
	Message string               `json:"message"`
}
*/

// A GetSessionResponse parameter model.
//
// Used as a response for getting user session.
//
// swagger:response getSessionResponse_v0
type GetGamesResponse_v0 struct {
	Game    presentationlayer.GamesView `json:"games"`
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
}
type GetGameResponse_v0 struct {
	Game    presentationlayer.GameView `json:"games"`
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
}
type EndSessionResponse_v0 struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// A CreateSessionRequest model.
//
// This is a CreateSessionRequest_v0 request model
//
// swagger:parameters handleCreateSession_v0
type CreateGameByPosRequest_v0 struct {
	Position mapmodels.LatLng `json:"latlng"`
}

type CreateGameByMapIDRequest_v0 struct {
	mapid string `json:"mid"`
}

// A CreateSessionResponse parameter model.
//
// Used as a response for create session request.
//
// swagger:response createSessionResponse_v0
type CreateGameResponse_v0 struct {
	Game    presentationlayer.GameView `json:"game"`
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
}

type BuyBuildingRequest_v0 struct {
	Buildingid string `json:"bid"`
	Spaceid    string `json:"sid"`
}

type BuyBuildingResponse_v0 struct {
	Game    presentationlayer.GameView `json:"game"`
	Message string                     `json:"message"`
	Success bool                       `json:"success"`
}

type GetEligibleBuildingsResponse_v0 struct {
	Buildings building.Buildings `json:"buildings"`
	Success   bool               `json:"success"`
	Message   string             `json:"message"`
}

type BuyUpgradeRequest_v0 struct {
	Upgradeid string `json:"uid"`
}

type BuyUpgradeResponse_v0 struct {
	Game    presentationlayer.GameView `json:"game"`
	Message string                     `json:"message"`
	Success bool                       `json:"success"`
}

type GetEligibleUpgradesResponse_v0 struct {
	Upgrades upgrade.Upgrades `json:"upgrades"`
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
}

// A ErrorResponse parameter model.
//
// Used as a response for errors.
//
// swagger:response errorResponse_v0
type ErrorResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

func (resource *Resource) DecodeRequestBody(w http.ResponseWriter, req *http.Request, target interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Request body parse error: %v", err.Error()))
		return err
	}
	return nil
}

func (resource *Resource) RenderError(w http.ResponseWriter, req *http.Request, status int, message string) {
	resource.Render(w, req, status, ErrorResponse_v0{
		Message: message,
		Success: false,
	})
}

func (resource *Resource) RenderUnauthorizedError(w http.ResponseWriter, req *http.Request, message string) {
	resource.Render(w, req, http.StatusUnauthorized, ErrorResponse_v0{
		Message: message,
		Success: false,
	})
}

// HandleCreateSession_v0 verify user's credentials and generates a JWT token if valid
// HandleCreateSession_v0 swagger:route POST /sessions sessions handleCreateSession_v0
//
// Creates a session for user
//
// Responses:
//    default: errorResponse_v0
//        200: createSessionResponse_v0
func (resource *Resource) HandleCreateGameByPos_v0(w http.ResponseWriter, req *http.Request) {
	gamerepo := resource.GameRepository(req)
	spacerepo := resource.SpaceRepository(req)
	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	var body CreateGameByPosRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	_game, _spaces, err := resource.CreateGameForPosition(body.Position)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	ok, err := gamerepo.GameExistsByUserVMapId(user.GetID(), _game.BuiltOnMapID.Hex())
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	if ok {
		resource.RenderError(w, req, http.StatusBadRequest, "User Already have a game close to that location")
		return
	}

	_game.BelongingUserID = bson.ObjectIdHex(user.GetID())
	newUserMapID := bson.NewObjectId()
	_game.Gm.GMap.ID = newUserMapID
	_game.GameCreationTime = time.Now()
	_game.LastIterationTime = time.Now()
	for idx, _ := range *_spaces {
		(*_spaces)[idx].BelongingMapID = newUserMapID

	}

	for _, f := range resource.ControllerHooks.PostCreateGameHook {
		f(resource, w, req, &PostCreateGameHookPayload{
			Ge: _game,
		})
	}

	err = gamerepo.CreateGame(_game)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	err = spacerepo.CreateSpaces(*_spaces)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	constructed, err := resource.attemptConstructGame_v0(*_game)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	gm := constructed.(*game.GameEnvironment)

	renderResource := resource.ViewRes(req).GetViewCreator()
	gv, err := renderResource.ConstructGameEnvView(*gm)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}
	resource.Render(w, req, http.StatusCreated, CreateGameResponse_v0{
		Game:    *gv,
		Success: true,
		Message: "Game Created",
	})
}

func (resource *Resource) HandleCreateGameByMapID_v0(w http.ResponseWriter, req *http.Request) {
	gamerepo := resource.GameRepository(req)
	spacerepo := resource.SpaceRepository(req)
	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	var body CreateGameByMapIDRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	ok, err := gamerepo.GameExistsByUserVMapId(user.GetID(), body.mapid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	if ok {
		resource.RenderError(w, req, http.StatusBadRequest, "User Already have a game close to that location")
		return
	}

	_game, _spaces, err := resource.CreateGameForMapID(body.mapid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	_game.BelongingUserID = bson.ObjectIdHex(user.GetID())
	newUserMapID := bson.NewObjectId()
	_game.Gm.GMap.ID = newUserMapID
	_game.GameCreationTime = time.Now()
	_game.LastIterationTime = time.Now()
	for idx, _ := range *_spaces {
		(*_spaces)[idx].BelongingMapID = newUserMapID

	}

	for _, f := range resource.ControllerHooks.PostCreateGameHook {
		f(resource, w, req, &PostCreateGameHookPayload{
			Ge: _game,
		})
	}

	err = gamerepo.CreateGame(_game)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	err = spacerepo.CreateSpaces(*_spaces)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	constructed, err := resource.attemptConstructGame_v0(*_game)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	gm := constructed.(*game.GameEnvironment)

	renderResource := resource.ViewRes(req).GetViewCreator()
	gv, err := renderResource.ConstructGameEnvView(*gm)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}
	resource.Render(w, req, http.StatusCreated, CreateGameResponse_v0{
		Game:    *gv,
		Success: true,
		Message: "Game Created",
	})
}

// HandleGetSession_v0 Get session details
// HandleGetSession_v0 swagger:route GET /sessions sessions handleGetSession_v0
//
// Gets the user session
//
// Responses:
//    default: errorResponse_v0
//        200: getSessionResponse_v0
func (resource *Resource) HandleGetGames_v0(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	userid := user.GetID()

	dbeds, err := resource.retrieveGames_v0("", userid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	if dbeds == nil {
		resource.RenderError(w, req, http.StatusBadRequest, "User has no games created")
		return
	}

	renderResource := resource.ViewRes(req).GetViewCreator()
	gvs, err := renderResource.ConstructGamesEnvView(dbeds.(game.GameEnvironments))
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}
	resource.Render(w, req, http.StatusOK, GetGamesResponse_v0{
		Game:    gvs,
		Success: true,
		Message: "Games retrieved",
	})
	return
}

func (resource *Resource) HandleGetGame_v0(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	params := mux.Vars(req)
	gameid := params["id"]
	userid := user.GetID()

	g, err := resource.retrieveGame_v0(gameid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	if g.GetBelongingUserID() != userid {
		resource.RenderError(w, req, http.StatusBadRequest, "Game does not belong to this user")
		return
	}

	//pretty.Println(renderResource.ConstructGameEnvView(*g.(*game.GameEnvironment)))
	renderResource := resource.ViewRes(req).GetViewCreator()

	gv, err := renderResource.ConstructGameEnvView(*g.(*game.GameEnvironment))
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}
	resource.Render(w, req, http.StatusOK, GetGameResponse_v0{
		Game:    *gv,
		Success: true,
		Message: "Game retrieved",
	})
	return
}

func (resource *Resource) HandleBuyUpgrade_v0(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	upgraderepo := resource.UpgradeRepository(req)
	resrepo := resource.ResourceRepository(req)

	user := users.GetUserCtx(ctx)

	params := mux.Vars(req)
	gameid := params["id"]

	upgradeid := params["uid"]

	g, err := resource.retrieveGame_v0(gameid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	if g.GetBelongingUserID() != user.GetID() {
		resource.RenderError(w, req, http.StatusBadRequest, "Game does not belong to this user")
		return
	}

	_u, err := upgraderepo.GetUpgradeById(upgradeid)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	//TODOSMT
	resourcesDB := resrepo.GetResources()
	upgrade, err := _u.FormUpgrade(resourcesDB)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	//fmt.Println(g)
	err = g.AddUpgrade(upgrade)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	//fmt.Println(upgrade)
	//g.AddUpgrade(u)
	//fmt.Println(g)
	err = resource.UpdateGame_v0(g)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	_g := g.(*game.GameEnvironment)

	renderResource := resource.ViewRes(req).GetViewCreator()
	gv, err := renderResource.ConstructGameEnvView(*_g)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}

	resource.Render(w, req, http.StatusOK, BuyUpgradeResponse_v0{
		Game:    *gv,
		Success: true,
		Message: "Upgrade Bought",
	})
}

func (resource *Resource) HandleGetEligibleUpgrades_v0(w http.ResponseWriter, req *http.Request) {
	//bMapper := resource.BuildingMapper(nil)
	ctx := req.Context()
	user := users.GetUserCtx(ctx)
	upgraderepo := resource.UpgradeRepository(req)
	resourcerepo := resource.ResourceRepository(req)

	params := mux.Vars(req)
	gameid := params["id"]

	userid := user.GetID()
	//CHANGE THIS ALGORITHM LATER
	_ = gameid
	_ = userid
	/*bs := bMapper.GetMappedBuildings()
	toRet := building.Buildings{}
	for _, v := range bs {
		toRet = append(toRet, *bMapper.MapBuildingDB(v))
	}
	*/
	_us := upgraderepo.GetUpgrades()
	_res := resourcerepo.GetResources()
	//fmt.Println(_us)
	us := upgrade.Upgrades{}
	for _, v := range _us {
		formed, err := v.FormUpgrade(_res)
		if err != nil {
			panic(err)
		}
		us = append(us, *formed)
	}
	resource.Render(w, req, http.StatusOK, GetEligibleUpgradesResponse_v0{
		Upgrades: us,
		Success:  true,
		Message:  "Upgrades Retrieved",
	})
}

func (resource *Resource) HandleBuyBuilding_v0(w http.ResponseWriter, req *http.Request) {
	bMapper := resource.BuildingMapper(req)
	//gamerepo := resource.GameRepository(req)

	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	params := mux.Vars(req)
	gameid := params["id"]

	if gameid == "" || !bson.IsObjectIdHex(gameid) {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Unknown game id"))
	}

	buildingid := params["bid"]
	spaceid := params["sid"]
	userid := user.GetID()

	_gms, err := resource.retrieveGames_v0(gameid, userid)
	gms := _gms.(game.GameEnvironments)
	if len(gms) != 1 {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Game not Found"))
	}
	retrievedGame := gms[0]

	if !retrievedGame.GetGame().GetGameMap().HasSpace(spaceid) {
		//fmt.Println(spaceid)
		//fmt.Println(game.GetGame().GetGameMap())
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Space not found"))
		return
	}

	//fmt.Println(retrievedGame.GetResources())

	buDB := bMapper.GetBuildingDBMapFor(buildingid)
	if buDB == nil {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Building not found"))
		return
	}

	bu := bMapper.MapBuildingDB(*buDB)
	if bu == nil {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Building not found"))
		return
	}
	created, err := retrievedGame.BuyBuilding(bu, spaceid, buildinginstance.FormFromBuildingOptions{ID: bson.NewObjectId().Hex(), Buildinglevel: 1})
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	created.SetBuiltTimeUnix(time.Now())
	created.SetLastProductionTimeUnix(time.Now())
	//fmt.Println(created.GetLastProductionTimeUnix)
	toadd, err := dbmodels.FormFromITileBuildableElement(created)
	toadd.BuildableVersion = buDB.Ver
	_, err = resource.HandleSetBuildingForSpace_v0(spaceid, *toadd)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	err = resource.UpdateGame_v0(&retrievedGame)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	renderResource := resource.ViewRes(req).GetViewCreator()
	gv, err := renderResource.ConstructGameEnvView(retrievedGame)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
	}
	resource.Render(w, req, http.StatusOK, BuyBuildingResponse_v0{
		Game:    *gv,
		Message: "Building Bought",
		Success: true,
	})
}

func (resource *Resource) HandleGetEligibleBuildings_v0(w http.ResponseWriter, req *http.Request) {
	bMapper := resource.BuildingMapper(nil)
	ctx := req.Context()
	user := users.GetUserCtx(ctx)

	params := mux.Vars(req)
	gameid := params["id"]

	userid := user.GetID()
	_ = gameid
	_ = userid
	//CHANGE THIS ALGORITHM LATER
	bs := bMapper.GetMappedBuildings()
	toRet := building.Buildings{}
	for _, v := range bs {
		toRet = append(toRet, *bMapper.MapBuildingDB(v))
	}

	//fmt.Println(toRet)
	resource.Render(w, req, http.StatusOK, GetEligibleBuildingsResponse_v0{
		Buildings: toRet,
		Success:   true,
		Message:   "Buildings Retrieved",
	})
}

func (resource *Resource) DeleteGameSessions_v0(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user := users.GetUserCtx(ctx)
	userid := user.GetID()

	//resource.GameManager.RemoveGamesByUserID(userid)
	_gs, ok := resource.GameManager.GetGamesByUserID(userid)
	if ok {
		gs := _gs.([]gamedomain.IGameEnvironment)
		for _, v := range gs {
			resource.GameManager.RemoveGameByGameID(v.GetGameID())
		}
	}

	resource.Render(w, req, http.StatusOK, EndSessionResponse_v0{
		Success: true,
		Message: "Session Ended",
	})

}

/*
func (resource *Resource) iterateGame_v0(gameid string) (gamedomain.IGameEnvironment, error) {
	resource.retrieveGame_v0(gameid)
}
*/
func (resource *Resource) retrieveGame_v0(gameid string) (gamedomain.IGameEnvironment, error) {
	cached, ok := resource.GameManager.GetGameByGameID(gameid)
	if ok {
		gm := cached.(gamedomain.IGameEnvironment)
		return gm, nil
	}
	repo := resource.GameRepository(nil)
	_game, err := repo.GetGameByGameId(gameid)
	constructed, err := resource.attemptConstructGame_v0(*_game)

	return constructed, err
}

func (resource *Resource) retrieveGames_v0(gameid string, userid string) (gamedomain.IGameEnvironments, error) {
	if gameid != "" {
		toRet := game.GameEnvironments{}
		cached, err := resource.retrieveGameFromGameCache_v0(gameid, userid)
		if err == nil {
			gm := cached.(*game.GameEnvironment)
			toRet = append(toRet, *gm)
			fmt.Println("game from cache")
			return toRet, nil
		}
	}
	fmt.Println("game from db")
	dbeds, err := resource.retrieveGamesFromDatabase_v0(gameid, userid)
	if err != nil {
		return nil, err
	}
	gms := dbeds.(game.GameEnvironments)
	//fmt.Println(len(gms))
	//for _, gm := range gms {
	for i := 0; i < len(gms); i++ {
		gm := gms[i]
		resource.GameManager.AddGame(&gm)
	}
	return dbeds, err
}

func (resource *Resource) retrieveGameFromGameCache_v0(gameid string, userid string) (gamedomain.IGameEnvironment, error) {
	cached, ok := resource.GameManager.GetGameByGameID(gameid)
	if ok {
		gm := cached.(gamedomain.IGameEnvironment)
		if gm.GetBelongingUserID() == userid {
			return gm, nil
		}
	}
	return nil, errors.New("Game Not In Cache Or User Not Allowed")
}

func (resource *Resource) retrieveGamesFromDatabase_v0(gameid string, userid string) (gamedomain.IGameEnvironments, error) {
	repo := resource.GameRepository(nil)
	_games, err := repo.GetGamesByUserId(userid)
	if err != nil {
		return nil, err
	}

	toRet := game.GameEnvironments{}
	if gameid != "" {
		for _, g := range _games {
			if g.GameID.Hex() == gameid {
				constructed, err := resource.attemptConstructGame_v0(g)
				if err != nil {
					return nil, errors.New("Game Data Corrupted")
				}
				toRet = append(toRet, *constructed.(*game.GameEnvironment))
				return toRet, nil
			}
		}
		return nil, errors.New("Game not found")
	}

	for _, _game := range _games {
		constructed, err := resource.attemptConstructGame_v0(_game)
		if err != nil {
			return nil, err
		}
		toRet = append(toRet, *constructed.(*game.GameEnvironment))
	}
	return toRet, nil
}

func (resource *Resource) attemptConstructGame_v0(g dbmodels.GameEnvironmentDB) (gamedomain.IGameEnvironment, error) {
	if !resource.validateGameMapIntegrity_v0(g) {
		return nil, errors.New("Map Integrity problem")
	}
	//upgrades := gm.Man.BoughtGameUpgradeIDs

	//builRepo := resource.BuildingRepository(req)
	resRepo := resource.ResourceRepository(nil)
	spaceRepo := resource.SpaceRepository(nil)
	bMapper := resource.BuildingMapper(nil)
	gmmap := g.Gm.GMap
	spaces := []gamemapdomain.ISpace{}
	dbspaces, err := spaceRepo.GetSpacesForGMap(gmmap.ID.Hex())
	if err != nil {
		return nil, errors.New("unknown spaces")
	}

	changes := datastructures.NewPQueueList(datastructures.MINPQ)
	//type todoFunc func(gamedomain.IGameEnvironment)

	//for _, v := range dbspaces {
	for i := 0; i < len(dbspaces); i++ {
		v := dbspaces[i]
		space := gamemap.NewSpace(v.ID.Hex(), v.InMapID, nil)
		placeholder := v.Element
		if placeholder != nil {
			//fmt.Println(v)
			buildable := placeholder.FormITileBuildableElement()
			curversion := placeholder.BuildableVersion
			/*if placeholder.BuildableType != "" {
			switch placeholder.BuildableType {
			default:
				fmt.Printf("unexpected type %T\n", placeholder) // %T prints whatever type t has
			case "buildinginstance":*/
			temp := buildable.(*buildinginstance.BuildingInstance)
			currentbu := bMapper.GetBuildingDBMapFor(temp.GetParentID())

			busDB := bMapper.GetBuildingsDBMapFor(currentbu.BID.Hex())

			var bu *building.Building = bMapper.MapBuildingDB(*currentbu)
			//for _, buDB := range busDB {
			for k := 0; k < len(busDB); k++ {
				buDB := busDB[k]
				comp, err := (&buDB.Ver).CompareVersion(&curversion)
				if err != nil {
					//return nil, err
					continue
				}

				if comp > 0 {
					changes.Push(func(g gamedomain.IGameEnvironment) {
						changebu := bMapper.MapBuildingDB(buDB)
						changebi := buildinginstance.FormFromBuilding(changebu, buildinginstance.FormFromBuildingOptions{ID: temp.GetUniqueID(), Buildinglevel: temp.GetLevel()})
						err := g.GetGame().GetGameMap().SetBuildingInstance(space.GetID(), changebi)
						if err == nil {
							changedbbuildable, err := dbmodels.FormFromITileBuildableElement(changebi)
							changedbbuildable.BuildableVersion = buDB.Ver

							if err == nil {
								spaceRepo.UpdateSpace(space.GetID(), &dbmodels.SpaceDB{Element: changedbbuildable})
							} else {
								fmt.Println(err)
							}
						} else {
							fmt.Println(err)
						}
					}, buDB.Ver.VerDateAfter.UnixNano())
				} else if comp == 0 {
					bu = bMapper.MapBuildingDB(buDB)
				}
			}
			//fmt.Println("here")
			//fmt.Println(pqueue)

			bi := buildinginstance.FormFromBuilding(bu, buildinginstance.FormFromBuildingOptions{ID: temp.GetUniqueID(), Buildinglevel: temp.GetLevel()})
			bi.SetBuiltTimeUnix(temp.GetBuiltTimeUnix())
			bi.SetLastProductionTimeUnix(temp.GetLastProductionTimeUnix())
			//fmt.Println(bi.GetLastProductionTimeUnix())
			space.AddElement(bi)
			//fmt.Println(space)

		}
		spaces = append(spaces, space)
	}

	for _, v := range spaces {
		for _, s := range dbspaces {
			if s.ID.Hex() == v.GetID() {
				for in, ns := range s.Around {
					for _, n := range ns {
						for _, v2 := range spaces {
							if v2.GetID() == n {
								v.AddNeighboorTo(v2, in)
							}
						}
					}
				}

			}
		}
	}

	Resources := generaldomain.IBProducts{}
	for _, v := range g.Resources {
		r, err := resRepo.GetResourceById(v.ID.Hex())
		if err != nil {
			return nil, errors.New("Resource Not found")
		}
		switch r.ResType {
		case "regularbproduct":
			Resources = append(Resources, general.NewRegularBProduct(v.ID.Hex(), v.ResourceCount))
		default:
			return nil, errors.New("Resource Type not recognized")
		}
	}

	GMap := gamemap.NewGameMap(spaces, gmmap.Version)
	gm := game.NewGame(GMap)
	toRet := game.NewGameEnvironment(g.GameID.Hex(), g.BelongingUserID.Hex(), g.BuiltOnMapID.Hex(), g.GameCreationTime, g.LastIterationTime, Resources, g.Man.BoughtGameUpgradeIDs, gm)

	// APPLYING VERSION UPDATES
	for todo, totick := changes.Pop(); todo != nil; todo, totick = changes.Pop() {
		//produceFor := time.Unix(totick, 0).UnixNano() - toRet.GetLastIterationTime().UnixNano()
		//produceDur := time.Duration(produceFor) * time.Nanosecond
		//toRet.TickFor(totick)
		err := resource.TickGame_v0(toRet, time.Unix(0, totick), toRet.GetLastIterationTime())
		if err != nil {
			panic(errors.New("Game couldnt tick"))
		}
		todoFc := todo.(func(g gamedomain.IGameEnvironment))
		todoFc(toRet)
		fmt.Println("done todo")
	}

	//produceFor := time.Now().UnixNano() - toRet.GetLastIterationTime().UnixNano()
	//produceDur := time.Duration(produceFor) * time.Nanosecond
	err = resource.TickGame_v0(toRet, time.Now(), toRet.GetLastIterationTime())
	if err != nil {
		panic(err)
	}
	//fmt.Println("Constructed Game:")
	//fmt.Println(toRet)
	return toRet, nil
}

func (resource *Resource) UpdateGame_v0(g gamedomain.IGameEnvironment) error {
	gamerepo := resource.GameRepository(nil)

	gameid := g.GetGameID()
	gdb, err := dbmodels.FormFromIGameEnvironment(g)
	if err != nil {
		return err
	}

	_, err = gamerepo.UpdateGame(gameid, gdb)

	return err
}

func (resource *Resource) TickGame_v0(g gamedomain.IGameEnvironment, t time.Time, u time.Time) error {
	produceFor := t.UnixNano() - u.UnixNano()
	produceDur := time.Duration(produceFor) * time.Nanosecond
	//gamerepo := resource.GameRepository(nil)
	spacerepo := resource.SpaceRepository(nil)
	//gameid := g.GetGameID()

	g.PlayFor(produceDur)

	err := resource.UpdateGame_v0(g)
	if err != nil {
		return err
	}

	g.GetGame().GetGameMap().ForSpaces(func(s gamemapdomain.ISpace) {
		if s.GetOccupier() == nil {
			return
		}

		buildable, err := dbmodels.FormFromITileBuildableElement(s.GetOccupier())
		if err != nil {
			panic(err)
			return
		}
		ns := dbmodels.SpaceDB{
			Element: buildable,
		}
		//fmt.Println(s.GetOccupier().(*buildinginstance.BuildingInstance).GetLastProductionTimeUnix())
		spacerepo.UpdateSpace(s.GetID(), &ns)
	})

	return nil
}

func (resource *Resource) HandleSetBuildingForSpace_v0(spaceid string, toadd dbmodels.TileBuildableElementDB) (*dbmodels.SpaceDB, error) {
	repo := resource.SpaceRepository(nil)
	upSpace := dbmodels.SpaceDB{
		Element: &toadd,
	}
	newSpace, err := repo.UpdateSpace(spaceid, &upSpace)
	if err != nil {
		return nil, err
	}
	return newSpace, nil
}

func (resource *Resource) validateGameMapIntegrity_v0(g dbmodels.GameEnvironmentDB) bool {
	gm := g.Gm.GMap
	refMapID := g.BuiltOnMapID

	mapRepo := resource.MapRepository(nil)
	_map, err := mapRepo.GetMapById(refMapID.Hex())
	if err != nil {
		fmt.Println(refMapID)
		fmt.Println(err)
		return false
	}
	curMap, ok := _map.(*maps.Map)
	if !ok {
		//WTF
	}
	//fmt.Println(gm.Version, curMap.MapVersion)
	c, err := (&gm.Version).CompareVersion(&curMap.MapVersion)
	if c != 0 {
		fmt.Println(gm.Version)
		fmt.Println(curMap.MapVersion)
		return false
	}
	return true
}
