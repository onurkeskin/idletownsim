package openCvExportWrappers

import (
	"bytes"
	"fmt"
)

type OpenCvReturn struct {
	Rels AllRelations
	Img  []byte
}

type AllRelations []EmptyPosRelations

type EmptyPosRelations struct {
	Self      EPosition          `json:"self" bson:"self"`
	SelfID    string             `json:"id,omitempty" bson:"_id,omitempty"`
	East      []EmptyPosRelation `json:"east" bson:"east"`
	South     []EmptyPosRelation `json:"south" bson:"south"`
	West      []EmptyPosRelation `json:"west" bson:"west"`
	North     []EmptyPosRelation `json:"north" bson:"north"`
	NorthWest []EmptyPosRelation `json:"northwest" bson:"northwest"`
	NorthEast []EmptyPosRelation `json:"northeast" bson:"northeast"`
	SouthWest []EmptyPosRelation `json:"southwest" bson:"southwest"`
	SouthEast []EmptyPosRelation `json:"southeast" bson:"southeast"`
}

func (p EmptyPosRelations) String() string {
	var buffer bytes.Buffer
	toAppendRect := fmt.Sprintf("Rect Pos:%s", p.Self.String())
	buffer.WriteString(toAppendRect)

	if len(p.East) > 0 {
		buffer.WriteString("\n")
		rightAppendPass := fmt.Sprintf("East Count: %d", len(p.East))
		buffer.WriteString(rightAppendPass)
	}
	for i := 0; i < len(p.East); i++ {
		buffer.WriteString("\n")
		toAppend := p.East[i].String()
		buffer.WriteString(toAppend)
	}

	if len(p.South) > 0 {
		buffer.WriteString("\n")
		bottomAppendPass := fmt.Sprintf("South Count: %d", len(p.South))
		buffer.WriteString(bottomAppendPass)
	}
	for i := 0; i < len(p.South); i++ {
		buffer.WriteString("\n")
		toAppend := p.South[i].String()
		buffer.WriteString(toAppend)
	}

	if len(p.West) > 0 {
		buffer.WriteString("\n")
		leftAppendPass := fmt.Sprintf("West Count: %d", len(p.West))
		buffer.WriteString(leftAppendPass)
	}
	for i := 0; i < len(p.West); i++ {
		buffer.WriteString("\n")
		toAppend := p.West[i].String()
		buffer.WriteString(toAppend)
	}

	if len(p.North) > 0 {
		buffer.WriteString("\n")
		topAppendPass := fmt.Sprintf("North Count: %d", len(p.North))
		buffer.WriteString(topAppendPass)
	}
	for i := 0; i < len(p.North); i++ {
		buffer.WriteString("\n")
		toAppend := p.North[i].String()
		buffer.WriteString(toAppend)
	}

	if len(p.NorthWest) > 0 {
		buffer.WriteString("\n")
		othersAppendPass := fmt.Sprintf("NorthWest Count: %d", len(p.NorthWest))
		buffer.WriteString(othersAppendPass)
	}
	for i := 0; i < len(p.NorthWest); i++ {
		buffer.WriteString("\n")
		toAppend := p.NorthWest[i].String()
		buffer.WriteString(toAppend)
	}

	if len(p.NorthEast) > 0 {
		buffer.WriteString("\n")
		othersAppendPass := fmt.Sprintf("NorthEast Count: %d", len(p.NorthEast))
		buffer.WriteString(othersAppendPass)
	}
	for i := 0; i < len(p.NorthEast); i++ {
		buffer.WriteString("\n")
		toAppend := p.NorthEast[i].String()
		buffer.WriteString(toAppend)
	}
	if len(p.SouthWest) > 0 {
		buffer.WriteString("\n")
		othersAppendPass := fmt.Sprintf("SouthWest Count: %d", len(p.SouthWest))
		buffer.WriteString(othersAppendPass)
	}
	for i := 0; i < len(p.SouthWest); i++ {
		buffer.WriteString("\n")
		toAppend := p.SouthWest[i].String()
		buffer.WriteString(toAppend)
	}
	if len(p.SouthEast) > 0 {
		buffer.WriteString("\n")
		othersAppendPass := fmt.Sprintf("SouthEast Count: %d", len(p.SouthEast))
		buffer.WriteString(othersAppendPass)
	}
	for i := 0; i < len(p.SouthEast); i++ {
		buffer.WriteString("\n")
		toAppend := p.SouthEast[i].String()
		buffer.WriteString(toAppend)
	}

	return buffer.String()
}

func (r AllRelations) String() string {
	var buffer bytes.Buffer
	for i := 0; i < len(r); i++ {
		toAdd := fmt.Sprintf("--------------Rect:%d-------------\n%s", i, r[i].String())
		buffer.WriteString(toAdd)
		if len(r)-1 != i {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
