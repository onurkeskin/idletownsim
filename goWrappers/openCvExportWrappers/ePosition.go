package openCvExportWrappers

import (
	"fmt"
)

type EPosition struct {
	P1, P2 BasicPoint
}

func (p EPosition) String() string {
	toRet := fmt.Sprintf("p1:%s,p2:%s", p.P1.String(), p.P2.String())
	return toRet
}
