package domain

type IEffectIssuer interface {
	GetDeployedEffects() []IEffect
	//GetUniqueID() string
}
