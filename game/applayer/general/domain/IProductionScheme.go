package domain

type IProductionScheme interface {
	GetProdType() string
	GetProdResult() IBProducts
}
