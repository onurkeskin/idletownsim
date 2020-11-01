package openCvExportWrappers

import (
	"fmt"
)

type EmptyPosRelation struct {
	FreePositionID string  `json:"freepositionid,omitempty" bson:"freepositionid,omitempty"`
	CenterAngle    float64 `json:"angle" bson:"angle"`
	CenterDistance int     `json:"distance" bson:"distance"`
}

func (p EmptyPosRelation) String() string {
	toRet := fmt.Sprintf("bound:%s ,angle:%f ,distance:%d", p.FreePositionID, p.CenterAngle, p.CenterDistance)
	return toRet
}
