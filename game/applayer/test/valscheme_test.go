package test

import (
	"github.com/app/game/applayer/general"
	"testing"
)

func TestMathScheme(t *testing.T) {
	sch, err := general.NewMathValScheme([]string{"+"}, []float64{2})
	if err {

	}

	var testNum float64 = 2
	var expected float64 = 4
	funcRes, _ := sch.CalculateValue(testNum)

	if funcRes != expected {
		t.Error("Error: Result:", funcRes, " Expected Value:", expected)
	}

	sch2, err2 := general.NewMathValScheme([]string{"*", "+"}, []float64{8, 2})
	if err2 {

	}
	var testNum2 float64 = 2
	var expected2 float64 = 18
	funcRes2, _ := sch2.CalculateValue(testNum2)

	if funcRes2 != expected2 {
		t.Error("Error: Result:", funcRes2, " Expected Value:", expected2)
	}

	_, err3 := general.NewMathValScheme([]string{"*", "+", "/"}, []float64{8, 3})
	if !err3 {
		t.Error("Error: Shouldnt be created")
	}

	t.Log("MathScheme Test Finished")
}
