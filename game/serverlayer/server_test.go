// The MIT License (MIT)

// Copyright (c) 2015 Hafiz Ismail

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/app/game/presentationlayer"
	"github.com/app/users"

	//"fmt"
	//"github.com/app/game/applayer/gamemap"
	//"github.com/app/game/applayer/general"
	//buildingdb "github.com/app/game/serverlayer/building"
	//"github.com/app/game/serverlayer/dbmodels"
	//usergame "github.com/app/game/serverlayer/usergame"
	//mapdb "github.com/app/map"
	//"github.com/app/server/middlewares/mongodb"
	"github.com/app/sessions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	//"time"
	//"errors"
	//resourcedb "github.com/app/game/serverlayer/gameresources"
	//spacedb "github.com/app/game/serverlayer/space"
	//upgradedb "github.com/app/game/serverlayer/upgrade"
	//mapmodels "github.com/app/map/mapmodels"
	//"github.com/app/server/middlewares/context"
	//"github.com/app/server/middlewares/renderer"
	//"github.com/app/server/server"
	//"github.com/app/users"
	//"io/ioutil"
)

/*
type MockGetGameResponse struct {
	Gs []GameResponse `json:"games"`
}
*/
type MockGetGameResponse struct {
	Gs presentationlayer.GamesView `json:"games,omitempty"`
}
type Gameenv struct {
	Gid string `json:"id,omitempty"`
	//	Uid string `json:"uid,omitempty"`
	//	Mid string `json:"mid,omitempty"`
	G Game `json:"game,omitempty"`
}
type Game struct {
	Gmap Gmap `json:"gmap,omitempty"`
}
type Gmap struct {
	Sids []ID `json:"spaces,omitempty"`
}

type ID struct {
	ID string `json:"id,omitempty"`
}

type MockMapResponse struct {
	Mapid   ID     `json:"map"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type MockBuildingsResponse struct {
	B []BuildingsJustId `json:"buildings,omitempty"`
}
type MockUpgradessResponse struct {
	U []UpgradesJustIds `json:"upgrades,omitempty"`
}

type BuildingsJustId struct {
	Bid string `json:"id,omitempty"`
}
type UpgradesJustIds struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	UpgradeID string `json:"uid,omitempty" bson:"uid,omitempty"`
}

var _ = Describe("Server", func() {
	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR TEST ONLY

	Describe("Server Tests", func() {
		Context("Integration Tests", func() {
			var (
				sessionToken string
				gameid       string
				mapid        string

				buyingbuildingid string
				spacebuildingid  string
				upgradebuyid     string
			)

			It("Shouldn't allow unauthorized access to /api/sessions", func() {
				// serve some urls
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/sessions", nil)
				s.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})

			It("Creating Account /api/users", func() {
				if Configuration.Create_user {
					recorderCreateUser := httptest.NewRecorder()
					var jsonStrCreateUser = []byte(`{"user": {"username": "onur", "password": "", "email": "onurkeskin@ku.edu.tr"}}`)
					request, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStrCreateUser))
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("Accept", "application/json")

					s.ServeHTTP(recorderCreateUser, request)
					var createuserresp users.CreateUserResponse_v0
					json.Unmarshal(recorderCreateUser.Body.Bytes(), &createuserresp)

					Expect(createuserresp.Success).To(Equal(true))
				}
			})

			It("Getting session token /api/sessions", func() {
				recorderCreateSession := httptest.NewRecorder()
				var jsonStrLogin = []byte(`{"username": "onur", "password": ""}`)
				request, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(jsonStrLogin))
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")

				s.ServeHTTP(recorderCreateSession, request)
				var createsessionresp sessions.CreateSessionResponse_v0
				json.Unmarshal(recorderCreateSession.Body.Bytes(), &createsessionresp)
				Expect(createsessionresp.Success).To(Equal(true))

				sessionToken = createsessionresp.Token
			})
			It("Creating Map /api/maps", func() {
				if Configuration.Create_game {
					recorderCreateMap := httptest.NewRecorder()
					var jsonStrCreateMap = []byte(`{"latlng": {"lat": 40.870458, "lng": 29.256617}}`)
					request, _ := http.NewRequest("POST", "/api/maps", bytes.NewBuffer(jsonStrCreateMap))
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("Accept", "application/json")
					request.Header.Add("Authorization", "Bearer "+sessionToken)
					s.ServeHTTP(recorderCreateMap, request)

					var mapresp MockMapResponse
					json.Unmarshal(recorderCreateMap.Body.Bytes(), &mapresp)
					Expect(mapresp.Success).To(Equal(true))

					mapid = mapresp.Mapid.ID
					fmt.Println(mapid)
				}
			})
			//-----------------------------------------------
			It("Creating Game /api/games/bypos", func() {
				if Configuration.Create_game {
					recorderCreateGame := httptest.NewRecorder()
					var jsonStrCreateGame = []byte(`{"latlng": {"lat": 40.870458, "lng": 29.256617}}`)
					request, _ := http.NewRequest("POST", "/api/game/bypos", bytes.NewBuffer(jsonStrCreateGame))
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("Accept", "application/json")
					request.Header.Add("Authorization", "Bearer "+sessionToken)
					s.ServeHTTP(recorderCreateGame, request)
					//fmt.Println(recorderCreateGame.Body)
				}
			})

			It("Creating Game /api/games/bymid", func() {
				if Configuration.Create_game {
					recorderCreateGame := httptest.NewRecorder()
					var jsonStrCreateGame = []byte(`{"mid": ""`)
					request, _ := http.NewRequest("POST", "/api/game/bymid", bytes.NewBuffer(jsonStrCreateGame))
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("Accept", "application/json")
					request.Header.Add("Authorization", "Bearer "+sessionToken)
					s.ServeHTTP(recorderCreateGame, request)
					//fmt.Println(recorderCreateGame.Body)
				}
			})

			It("Getting all maps /api/maps/?lat=40.870458&lng=29.256617"+upgradebuyid, func() {
				recorderGetMaps := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/maps?lat=40.870458&lng=29.256617", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetMaps, request)
				//fmt.Println(recorderGetMaps.Body)
			})
			//-----------------------------------------------

			It("Getting game /api/games", func() {
				recorderGetGame := httptest.NewRecorder()
				//var jsonStrGameId = []byte(`{"gid": "` + gameid + `"}`)
				request, _ := http.NewRequest("GET", "/api/games", nil /*bytes.NewBuffer(jsonStrGameId)*/)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetGame, request)
				var getgameresp MockGetGameResponse //presentationlayer.GamesView
				json.Unmarshal(recorderGetGame.Body.Bytes(), &getgameresp)
				//fmt.Println(recorderGetGame.Body)
				//fmt.Println(getgameresp)
				gameid = getgameresp.Gs[0].GameID
				spacebuildingid = getgameresp.Gs[0].Spaces[0].SpaceID
				//fmt.Println(gameid)
				//fmt.Println(spacebuildingid)
				//fmt.Println(recorderGetGame.Body)
			})

			It("Getting game /api/game", func() {
				recorderGetGame := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/game/"+gameid, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetGame, request)
				//fmt.Println(recorderGetGame.Body)
			})

			It("Deleting game session /api/games", func() {
				recorderGetGame := httptest.NewRecorder()
				request, _ := http.NewRequest("DELETE", "/api/games", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetGame, request)
				//fmt.Println(recorderGetGame.Body)
			})

			It("Getting buildings for game /api/game/"+gameid+"/buildings", func() {
				recorderGetBuildings := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/game/"+gameid+"/buildings", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetBuildings, request)
				var buildingsresp MockBuildingsResponse
				json.Unmarshal(recorderGetBuildings.Body.Bytes(), &buildingsresp)
				buyingbuildingid = buildingsresp.B[0].Bid
				//fmt.Println(buyingbuildingid)
			})

			It("Buying building for game /api/game/"+gameid+"/space/"+spacebuildingid+"/building/buy/"+buyingbuildingid, func() {
				recorderBuyBuilding := httptest.NewRecorder()
				request, _ := http.NewRequest("POST", "/api/game/"+gameid+"/space/"+spacebuildingid+"/building/buy/"+buyingbuildingid, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderBuyBuilding, request)
				//fmt.Println(recorderBuyBuilding.Body)
			})

			It("Getting upgrades for game /api/game/"+gameid+"/upgrades", func() {
				recorderGetUpgrades := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/game/"+gameid+"/upgrades", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderGetUpgrades, request)
				var upgradesresp MockUpgradessResponse
				json.Unmarshal(recorderGetUpgrades.Body.Bytes(), &upgradesresp)
				//fmt.Println(recorderGetUpgrades.Body)
				upgradebuyid = upgradesresp.U[0].ID
				//fmt.Println(upgradesresp)
			})

			It("Buying upgrade for game /api/game/"+gameid+"/upgrade/buy/"+upgradebuyid, func() {
				recorderBuyUpgrade := httptest.NewRecorder()
				request, _ := http.NewRequest("POST", "/api/game/"+gameid+"/upgrade/buy/"+upgradebuyid, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("Accept", "application/json")
				request.Header.Add("Authorization", "Bearer "+sessionToken)
				s.ServeHTTP(recorderBuyUpgrade, request)
				//fmt.Println(recorderBuyUpgrade.Body)
			})

		})
	})
})
