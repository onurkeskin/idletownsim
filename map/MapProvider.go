package maps

import (
	"bytes"
	//b64 "encoding/base64"
	///"fmt"
	"github.com/app/appenginehelpers"
	"github.com/app/goWrappers/openCvExportWrappers"
	mapdomain "github.com/app/map/domain"
	googlemaps "googlemaps.github.io/maps"
	"strings"
	//"github.com/kr/pretty"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	//"strconv"
	"github.com/app/map/mapmodels"
	"time"
)

const (
	requestedMapStyle string = "&style=feature:all|color:0x000000&style=feature:road|color:0xFFFFFF&style=feature:road.local|color:0xFFFFFF&zoom=18&size=400x400&style=feature:all|element:labels|visibility:off"
)

type MapProvider struct {
	Retriever MapImageProvider
}

func (m *MapProvider) GetMapByLatLang(latlng mapmodels.LatLng) (mapdomain.IMap, error) {
	res, err := MapImageByProvider(m.Retriever, latlng)
	if err != nil {
		return nil, err
	}

	parsed := m.ParseMapImage(res.Image)
	toRet := &Map{
		MapIdentifierAdress: res.MapIdentifierAdress,
		MapCompleteAddress:  res.MapCompleteAddress,
		MapRaw:              parsed.Img,
		MapFundamentalCoordinates: latlng,
		ParsedRelations:           parsed.Rels,

		Status:           "created",
		LastModifiedDate: time.Now(),
		CreatedDate:      time.Now(),
	}
	return toRet, nil
}

func MapImageByGoogleMaps(latlng mapmodels.LatLng) (v MapProviderReturn, err error) {
	prov := appengessentials.GetInstance()
	mapKey := appengessentials.GetKeyFromFile("googlemapsStaticKey", "keys/googlestaticmapskey")
	// MAKE SOME KEY LOADER IN FUTURE
	prov.SetKeyFor(mapKey[0], mapKey[1])
	//fmt.Println(prov.GetKeyFor("googlemapsStaticKey"))
	c, err := googlemaps.NewClient(googlemaps.WithAPIKey(prov.GetKeyFor("googlemapsStaticKey")))
	//c, err := googlemaps.NewClient(googlemaps.WithAPIKey("AIzaSyA6h41WJt6SMhjESlMa2JIw0yMN4CIa4uI"))
	if err != nil {
		panic(err)
	}
	userLatLng := googlemaps.LatLng{Lat: latlng.Lat, Lng: latlng.Lng}

	r := &googlemaps.GeocodingRequest{
		LatLng: &userLatLng,
	}
	resp, err := c.ReverseGeocode(context.Background(), r)
	if err != nil {
		panic(err)
	}

	overallLocStr := LocationParserSpe(resp[0].AddressComponents)

	res, err := GetMapImage(overallLocStr)

	toRet := MapProviderReturn{
		Image:               res,
		MapIdentifierAdress: overallLocStr,
		MapCompleteAddress:  overallLocStr,
	}
	return toRet, err
}

type MapProviderReturn struct {
	Image               []byte
	MapIdentifierAdress string
	MapCompleteAddress  string
}

func MapImageByProvider(provider MapImageProvider, latlng mapmodels.LatLng) (v MapProviderReturn, err error) {
	return provider(latlng)
}

type MapImageProvider func(mapmodels.LatLng) (v MapProviderReturn, err error)

func GetMapImage(placeName string) (image []byte, err error) {
	googlemapsKey := appengessentials.GetInstance().GetKeyFor("googlemapsStaticKey")

	var buffer bytes.Buffer
	buffer.WriteString("https://maps.googleapis.com/maps/api/staticmap?center=")
	buffer.WriteString(placeName)
	buffer.WriteString(requestedMapStyle)
	buffer.WriteString("&key=")
	buffer.WriteString(googlemapsKey)
	//fmt.Println(buffer.String())

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get(buffer.String())
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.Body != nil {
		image, err = ioutil.ReadAll(response.Body)
	}

	//pretty.Println(b64.StdEncoding.EncodeToString(image))
	return image, err
}

func (m *MapProvider) ParseMapImage(image []byte) openCvExportWrappers.OpenCvReturn {
	test := openCvExportWrappers.DoStuff(image)
	return test
}

func WriteImageToFile(img []byte, name string, ex string) {
	err := ioutil.WriteFile(name+"."+ex, img, 0644)
	if err != nil {
		panic(err)
	}
}

func LocationParserSpe(address []googlemaps.AddressComponent) string {
	var toRet string

	var eligible []googlemaps.AddressComponent
	for _, el := range address {
		if checkImportantAddressParam(el.Types) {
			eligible = append(eligible, el)
			//fmt.Printf("%v %v\n", "eligible: "+el.LongName+" type: ", el.Types)
		} else {
			//fmt.Printf("%v %v\n", "ineligible: "+el.LongName+" type: ", el.Types)
		}
	}

	for i := 0; i < len(eligible); i++ {
		toRet += swapSpacesWithPlus(eligible[i].LongName)
		toRet += ","
	}

	if last := len(toRet) - 1; last >= 0 && toRet[last] == ',' {
		toRet = toRet[:last]
	}

	return toRet
}

func checkImportantAddressParam(param []string) bool {
	for _, el := range param {
		switch el {
		case "administrative_area_level_1":
			fallthrough
		case "administrative_area_level_2":
			fallthrough
		case "administrative_area_level_3":
			fallthrough
		case "administrative_area_level_4":
			fallthrough
		case "country":
			fallthrough
		case "route":
			return true
		}

	}
	return false
}

func swapSpacesWithPlus(param string) string {
	return strings.Replace(param, " ", "+", -1)
}
