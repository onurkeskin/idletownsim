package domain

type IValScheme interface {
	CalculateValue(i interface{}) (interface{}, bool)
}
