package general

import (
	"fmt"
	"math"
)

type MathValScheme struct {
	Operators []string  `json:"operators,omitempty"" bson:"operators,omitempty""`
	Numbers   []float64 `json:"numbers,omitempty"" bson:"numbers,omitempty""`
}

func NewMathValScheme(Operators []string, Numbers []float64) (*MathValScheme, bool) {
	if len(Operators) != len(Numbers) {
		return nil, true
	}

	return &MathValScheme{Operators, Numbers}, false
}

func (m *MathValScheme) CalculateValue(i interface{}) (interface{}, bool) {
	a := i.(float64)
	for k := 0; k < len(m.Operators); k++ {
		doOperation(&a, m.Operators[k], m.Numbers[k])
	}

	return a, false
}

func doOperation(num *float64, operator string, operand float64) bool {
	switch operator {
	case "+":
		*num = *num + operand
		return true
	case "-":
		*num = *num - operand
		return true
	case "*":
		*num = *num * operand
		return true
	case "/":
		*num = *num / operand
		return true
	case "%":
		*num = math.Mod(*num, operand)
		return true
	}
	return false
}

func (g *MathValScheme) String() string {
	str := ""
	str += fmt.Sprintf("Math Eq:")
	for i := 0; i < len(g.Operators); i++ {
		str += fmt.Sprintf("(")
	}
	str += fmt.Sprintf("X")
	for i := 0; i < len(g.Operators); i++ {
		str += fmt.Sprintf("%s%f)", g.Operators[i], g.Numbers[i])
	}
	return str
}
