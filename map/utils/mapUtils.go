package maputils

import (
	"bytes"
	"github.com/app/appenginehelpers"
	"github.com/lazywei/go-opencv/opencv"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetMapImage(latitude float64, longitude float64) (image []byte, err error) {
	mapsKey := appengessentials.GetInstance().GetKeyFor("MapsStaticKey")

	var buffer bytes.Buffer
	buffer.WriteString("https://maps.googleapis.com/maps/api/staticmap?center=")
	buffer.WriteString(strconv.FormatFloat(latitude, 'f', 6, 64))
	buffer.WriteString(",")
	buffer.WriteString(strconv.FormatFloat(longitude, 'f', 6, 64))
	buffer.WriteString("&zoom=12&size=400x400&key=")
	buffer.WriteString(mapsKey)

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, _ := netClient.Get(buffer.String())

	if response.Body != nil {
		image, err = ioutil.ReadAll(response.Body)
	}

	return image, err
}

func ParseMapImage(image []byte) {
	image := opencv.DecodeImageMem(image)
	var win opencv.Window
	win.ShowImage(image)

	opencv.WaitKey(0)
}
