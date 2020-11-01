package sessions

import (
	"time"
)

type ITokenClaims interface {
	GetJTI() string
	GetExpireAt() time.Time
}

type IToken interface {
	IsValid() bool
}

type IRevokedToken interface {
}
