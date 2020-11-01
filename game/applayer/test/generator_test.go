package test

import (
	general "github.com/app/game/applayer/general"
	generator "github.com/app/game/applayer/generator"
	"testing"
)

func TestGenerator(t *testing.T) {
	typ1 := "a"
	typ2 := "b"
	typ1Inc := float64(5)
	typ2Dec := float64(-5)

	increase := general.NewRegularBProduct(typ1, typ1Inc)
	//incArr := domain.IBProducts{increase}

	decrease := general.NewRegularBProduct(typ2, typ2Dec)
	//decArr := domain.IBProducts{decrease}

	b := generator.NewBasicGenerator()
	b.AddIncrease(increase)
	b.AddIncrease(decrease)

	med := b.Generate(nil)
	if len(med) == 0 {
		t.Error("Empty after non null production")
	}

	ExpectedRes1 := float64(5)
	ExpectedRes2 := float64(-5)
	for _, check := range med {
		if check.GetType() == typ1 {
			if check.GetValue() != ExpectedRes1 {
				t.Log("Current type1 :", typ1, " val:", check.GetValue(), " Expected: ", ExpectedRes1)
				t.Error("Wrong production")
			}
		} else if check.GetType() == typ2 {
			if check.GetValue() != ExpectedRes2 {
				t.Log("Current type1 :", typ2, " val:", check.GetValue(), " Expected: ", ExpectedRes2)
				t.Error("Wrong production")
			}
		} else {
			t.Log("Wrong Production type")
			t.Error("Wrong production")
		}
	}

}
