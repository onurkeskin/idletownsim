package version

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	Ver          string    `json:"version" bson:"version"`
	VerDateAfter time.Time `json:"validafter" bson:"validafter"`
}

func (d *Version) MarshalJSON() ([]byte, error) {
	type toJsonVer Version
	return json.Marshal(&struct {
		*toJsonVer
		VerDateAfter string `json:"validafter" bson:"validafter"`
	}{
		toJsonVer:    (*toJsonVer)(d),
		VerDateAfter: d.VerDateAfter.Format("2006-01-02T15:04:05.999999Z"),
	})
}

func (v *Version) CompareVersion(comp *Version) (int, error) {
	var validVersionString = regexp.MustCompile(`^(\d+)((\.{1}\d+)*)(\.{0})$`)
	selfok := validVersionString.MatchString(v.Ver)
	if !selfok {
		return 0, errors.New("Self version is not eligible to be a version")
	}
	compok := validVersionString.MatchString(comp.Ver)
	if !compok {
		return 0, errors.New("Compared version is not eligible to be a version")
	}

	vNums := strings.Split(v.Ver, ".")
	compNums := strings.Split(comp.Ver, ".")
	for i, _ := range vNums {
		self, err := strconv.Atoi(vNums[i])
		if err != nil {
			return 0, err
		}
		other, err := strconv.Atoi(compNums[i])
		if err != nil {
			return 0, err
		}
		if self > other {
			return 1, nil
		} else if self < other {
			return -1, nil
		}
	}
	if len(vNums) > len(compNums) {
		return 1, nil
	} else if len(vNums) < len(compNums) {
		return -1, nil
	}

	return 0, nil
}
