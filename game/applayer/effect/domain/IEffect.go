package domain

type IEffect interface {
	GetUniqueID() string

	ApplyEffectGlobal(anything interface{})
	RemoveEffect()
	ReapplyEffect()

	GetIssuer() IEffectIssuer
	SetIssuer(iss IEffectIssuer)
}
