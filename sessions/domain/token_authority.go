package sessions

import ()

type ITokenAuthority interface {
	CreateNewSessionToken(claims ITokenClaims) (string, error)
	VerifyTokenString(tokenStr string) (IToken, ITokenClaims, error)
}
