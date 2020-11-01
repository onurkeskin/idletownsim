package test

import (
	"github.com/app/game/applayer/general"
	"github.com/app/game/applayer/general/domain"

	"testing"
)

func TestObjectProperties(t *testing.T) {
	p1 := general.NewObjectProperties()
	p2 := general.NewObjectProperties()

	p1.AddProperty("n1", "v1")
	p1.AddProperty("n2", "v2")
	p1.AddProperty("n3", "v3")

	p2.AddProperty("n1", "v1")
	p2.AddProperty("n2", "v2")
	p2.AddProperty("n3", "v3")

	res := p1.SatisfiesType(p2, domain.CheckTypeExact)
	if !res {
		t.Error("Wrong")
	}

	p2.RemoveProperty("n3", "v3")
	res = p1.SatisfiesType(p2, domain.CheckTypeExact)
	if res {
		t.Error("Wrong")
	}

	res2, _ := p2.GetProperty("n1")
	if res2[0] != "v1" {
		t.Error("Wrong")
	}

}
