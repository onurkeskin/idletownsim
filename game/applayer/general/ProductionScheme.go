package general

type ProductionScheme struct {
	productionType   string
	productionResult interface{}
}

func (pScheme ProductionScheme) GetProdType() string {
	return pScheme.productionType
}

func (pScheme ProductionScheme) GetProdResult() interface{} {
	return pScheme.productionResult
}
