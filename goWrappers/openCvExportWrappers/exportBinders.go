package openCvExportWrappers

/*
#cgo CFLAGS: -I${SRCDIR}/../../olib -I/usr/include
#cgo LDFLAGS: -L${SRCDIR}/../../olib -Wl,-rpath,${SRCDIR}/../../olib -lmapProvider
#cgo pkg-config: opencv
#include <genericHeader.hpp>
#include "cWrapper.h"
#include <stdlib.h>
*/
import "C"
import (
	//"fmt"
	"strconv"
	"unsafe"
)

func DoStuff(image []byte) OpenCvReturn {
	/*
		cdata := C.malloc(C.size_t(len(image)))
		defer C.free(cdata)

		copy((*[1 << 24]byte)(cdata)[0:len(image)], image)
		possUnsafePtr := C.basic((*C.uchar)(cdata), C.int(len(image)))
		toRet := ParseIntoGo(possUnsafePtr)
		defer FreeCExports(possUnsafePtr)
		return toRet
	*/
	p := C.malloc(C.size_t(len(image)))
	defer C.free(p)
	cBuf := (*[1 << 30]byte)(p)
	copy(cBuf[:], image)
	possUnsafePtr := C.basic((*C.uchar)(p), C.int(len(image)))
	toRet := ParseIntoGo(possUnsafePtr)
	defer FreeCExports(possUnsafePtr)
	return toRet
}

func ParseMeAEmptyPosRel(r C.cEPosition) EPosition {
	var parentp1x int = int(r.p1.xval)
	var parentp1y int = int(r.p1.yval)
	var parentp2x int = int(r.p2.xval)
	var parentp2y int = int(r.p2.yval)
	p1B := BasicPoint{parentp1x, parentp1y}
	p2B := BasicPoint{parentp2x, parentp2y}
	pEP := EPosition{p1B, p2B}

	return pEP
}

func FreeCExports(possUnsafePtr C.cExports) {
	elements := possUnsafePtr.exports
	size := int((possUnsafePtr).count)
	//img := possUnsafePtr.modimg
	var teamSlice []C.cEmptyPosRelations = (*[1 << 30]C.cEmptyPosRelations)(unsafe.Pointer(elements))[:size:size]
	for i := 0; i < size; i++ {
		C.free(unsafe.Pointer(teamSlice[i].east))
		C.free(unsafe.Pointer(teamSlice[i].south))
		C.free(unsafe.Pointer(teamSlice[i].west))
		C.free(unsafe.Pointer(teamSlice[i].north))
		C.free(unsafe.Pointer(teamSlice[i].southeast))
		C.free(unsafe.Pointer(teamSlice[i].southwest))
		C.free(unsafe.Pointer(teamSlice[i].northeast))
		C.free(unsafe.Pointer(teamSlice[i].northwest))
	}
	C.free(unsafe.Pointer(elements))
	// No need to free because the image is not allocated by a malloc handled by cpp probably
	//C.free(unsafe.Pointer(possUnsafePtr.modimg))
}

func CompareEmptyRelParents(crel C.cEPosition, rel EPosition) bool {
	toCmp := ParseMeAEmptyPosRel(crel)
	if toCmp == rel {
		return true
	}

	return false
}

func ParseIntoGo(possUnsafePtr C.cExports) OpenCvReturn {
	elements := possUnsafePtr.exports
	size := int((possUnsafePtr).count)
	//imgsize := int(possUnsafePtr.imgsize)
	img := C.GoBytes(unsafe.Pointer(possUnsafePtr.modimg), possUnsafePtr.imgsize)
	//img := (*[1 << 30]C.uchar)(unsafe.Pointer(possUnsafePtr.modimg))[:imgsize:imgsize]
	//fmt.Printf("Size = %v \n", imgsize)
	//fmt.Println(img)
	rels := make([]EmptyPosRelations, size)

	var teamSlice []C.cEmptyPosRelations = (*[1 << 30]C.cEmptyPosRelations)(unsafe.Pointer(elements))[:size:size]
	for i := 0; i < size; i++ {
		pEP := ParseMeAEmptyPosRel(teamSlice[i].parent)

		curParentRel := EmptyPosRelations{
			Self:   pEP,
			SelfID: strconv.Itoa(i),
		}

		rels[i] = curParentRel
	}

	for x := 0; x < size; x++ {
		//r := rels[x]
		for i := 0; i < size; i++ {
			if !CompareEmptyRelParents(teamSlice[i].parent, rels[x].Self) {
				continue
			}

			curR := teamSlice[i].east
			curSlice := (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].eastSize:teamSlice[i].eastSize]
			for rN := 0; rN < int(teamSlice[i].eastSize); rN++ {
				curS := curSlice[rN]
				//pEP := ParseMeAEmptyPosRel(curS.freePosition)
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].East = append(rels[x].East, rel)
				}
			}

			curR = teamSlice[i].south
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].southSize:teamSlice[i].southSize]
			for rN := 0; rN < int(teamSlice[i].southSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].South = append(rels[x].South, rel)
				}
			}

			curR = teamSlice[i].west
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].westSize:teamSlice[i].westSize]
			for rN := 0; rN < int(teamSlice[i].westSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].West = append(rels[x].West, rel)
				}
			}

			curR = teamSlice[i].north
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].northSize:teamSlice[i].northSize]
			for rN := 0; rN < int(teamSlice[i].northSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].North = append(rels[x].North, rel)
				}
			}

			curR = teamSlice[i].southeast
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].southeastSize:teamSlice[i].southeastSize]
			for rN := 0; rN < int(teamSlice[i].southeastSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].SouthEast = append(rels[x].SouthEast, rel)
				}
			}
			curR = teamSlice[i].southwest
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].southwestSize:teamSlice[i].southwestSize]
			for rN := 0; rN < int(teamSlice[i].southwestSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].SouthWest = append(rels[x].SouthWest, rel)
				}
			}
			curR = teamSlice[i].northeast
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].northeastSize:teamSlice[i].northeastSize]
			for rN := 0; rN < int(teamSlice[i].northeastSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].NorthEast = append(rels[x].NorthEast, rel)
				}
			}
			curR = teamSlice[i].northwest
			curSlice = (*[1 << 30]C.cEmptyPosRelation)(unsafe.Pointer(curR))[:teamSlice[i].northwestSize:teamSlice[i].northwestSize]
			for rN := 0; rN < int(teamSlice[i].northwestSize); rN++ {
				curS := curSlice[rN]
				for _, rt := range rels {
					if !CompareEmptyRelParents(curS.freePosition, rt.Self) {
						continue
					}
					var rel EmptyPosRelation = EmptyPosRelation{rt.SelfID, float64(curS.centerAngle), int(curS.centerDistance)}
					rels[x].NorthWest = append(rels[x].NorthWest, rel)
				}
			}
		}
	}

	toRet := OpenCvReturn{
		Img:  img,
		Rels: rels,
	}

	return toRet
}
