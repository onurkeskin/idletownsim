package domain

type CheckType int

const (
	CheckTypeAll CheckType = iota
	CheckTypeAny
	CheckTypeExact
)

func (c CheckType) String() string {
	switch c {
	case CheckTypeAll:
		return "All"
	case CheckTypeAny:
		return "Any"
	case CheckTypeExact:
		return "Exact"
	default:
		return "Unknown"
	}
}

type IObjectProperty interface {
	GetParamName() string
	GetParamValue() string
}

type IObjectProperties interface {
	Count() int
	PropertiesSlice() []IObjectProperty
	GetProperty(paramName string) ([]string, error)
	AddProperty(paramName, paramValue string) error
	RemoveProperty(paramName, paramValue string) error
	SatisfiesType(check IObjectProperties, chkType CheckType) bool
}
