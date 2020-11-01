package openCvExportWrappers

import (
	"fmt"
)

type BasicPoint struct {
	Xval, Yval int
}

func (p BasicPoint) String() string {
	toRet := fmt.Sprintf("X:%d Y:%d", p.Xval, p.Yval)
	return toRet
}
