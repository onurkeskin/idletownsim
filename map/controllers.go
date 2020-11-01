package maps

import (
	"encoding/json"
	"fmt"
	mapdomain "github.com/app/map/domain"
	"github.com/app/map/mapmodels"
	"strconv"
	//"net/http/httputil"
	//"fmt"
	//"github.com/app/server/domain"
	"github.com/gorilla/mux"
	//"gopkg.in/mgo.v2/bson"
	//"log"
	"net/http"
)

type GetMapResponse_v0 struct {
	Map     mapdomain.IMap `json:"map"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
}

type CreateMapRequest_v0 struct {
	Position mapmodels.LatLng `json:"latlng"`
}

type CreateMapResponse_v0 struct {
	Map     mapdomain.IMap `json:"map"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
}

type GetCloseMapsResponse_v0 struct {
	Maps    mapdomain.IMaps `json:"maps"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
}

// A ErrorResponse parameter model.
//
// Used as a response for errors.
//
// swagger:response errorResponse_v0
type ErrorResponse_v0 struct {
	// in: body
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

func (resource *Resource) DecodeRequestBody(w http.ResponseWriter, req *http.Request, target interface{}) error {
	/*requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	*/
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

func (resource *Resource) HandleGetMap_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	repo := resource.MapRepository(req)
	_maps, err := repo.GetMapById(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Couldn't find map")
		return
	}
	retmap := _maps.(*Map)

	resource.Render(w, req, http.StatusOK, GetMapResponse_v0{
		Map:     retmap,
		Message: "Map retrieved",
		Success: true,
	})
}

func (resource *Resource) HandleGetCloseMaps_v0(w http.ResponseWriter, req *http.Request) {
	queryparams := req.URL.Query()
	_lat := queryparams.Get("lat")
	_lng := queryparams.Get("lng")
	lat, err := strconv.ParseFloat(_lat, 64)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Invalid Latitude")
		return
	}
	lng, err := strconv.ParseFloat(_lng, 64)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Invalid Longitude")
		return
	}

	pos := mapmodels.LatLng{
		Lat: lat,
		Lng: lng,
	}
	repo := resource.MapRepository(req)

	maps, err := GetClosestMapsByPosition(repo, pos, 1)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Couldn't find map")
		return
	}

	resource.Render(w, req, http.StatusOK, GetCloseMapsResponse_v0{
		Maps:    maps,
		Message: "Map retrieved",
		Success: true,
	})
}

func (resource *Resource) HandleCreateMapLatLng_v0(w http.ResponseWriter, req *http.Request) {
	repo := resource.MapRepository(req)

	var body CreateMapRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	m, err := CreateMapForPosition(resource.MapProvider, repo, body.Position)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	resource.Render(w, req, http.StatusOK, CreateMapResponse_v0{
		Map:     m,
		Message: "Map Created",
		Success: true,
	})
}
