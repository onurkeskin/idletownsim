package effect

import (
	"fmt"
	"github.com/app/game/applayer/general"
)

type GenerationEffect struct {
	Effect          `json:"effectbase,omitempty" bson:"effectbase,omitempty"`
	Scheme          *general.MathValScheme `json:"effectscheme" bson:"effectscheme"`
	TargetResources []string               `json:"effecttargetresource" bson:"effecttargetresource"`
}

func (g *GenerationEffect) String() string {
	str := ""
	str += fmt.Sprintf("Effect:%s", g.Effect.String())
	str += fmt.Sprintf("scheme:%s ", g.Scheme)
	str += fmt.Sprintf("targetResources:%v", g.TargetResources)
	return str
}
