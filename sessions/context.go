package sessions

import (
	. "github.com/app/sessions/domain"
	"golang.org/x/net/context"
	//"strings"

	"github.com/app/server/domain"
	//"net/http"
)

const TokenAuthorityKey domain.ContextKey = "sessionTokenAuthorityKey"
const TokenClaimsKey domain.ContextKey = "sessionTokenClaimsKey"

func GetTokenAuthorityCtx(ctx context.Context) ITokenAuthority {
	return ctx.Value(TokenAuthorityKey).(ITokenAuthority)
}

func GetAuthenticatedClaimsCtx(ctx context.Context) ITokenClaims {
	return ctx.Value(TokenClaimsKey).(ITokenClaims)
}

/*
func newContextWithTokenAuthority(ctx context.Context, req *http.Request) context.Context {
	authHeaderString := req.Header.Get("Authorization")
	if authHeaderString != "" {
		tokens := strings.Split(authHeaderString, " ")
		if len(tokens) != 2 || (len(tokens) > 0 && strings.ToUpper(tokens[0]) != "BEARER") {
			//resource.RenderUnauthorizedError(w, req, "Invalid format, expected Authorization: Bearer {token}")
			return ctx
		}

		tokenString := tokens[1]
		t, c, err := resource.TokenAuthority.VerifyTokenString(tokenString)
		if err != nil {
			//resource.RenderUnauthorizedError(w, req, "Unable to verify token string")
			return ctx
		}
		token := t.(*Token)
		claims := c.(*TokenClaims)
	}

	return context.WithValue(ctx, TokenAuthorityKey, Token)
}

func newContextWithAuthenticatedClaims(ctx context.Context, req *http.Request) context.Context {
	authHeaderString := req.Header.Get("Authorization")
	if authHeaderString != "" {
		tokens := strings.Split(authHeaderString, " ")
		if len(tokens) != 2 || (len(tokens) > 0 && strings.ToUpper(tokens[0]) != "BEARER") {
			//resource.RenderUnauthorizedError(w, req, "Invalid format, expected Authorization: Bearer {token}")
			return ctx
		}

		tokenString := tokens[1]
		_, c, err := resource.TokenAuthority.VerifyTokenString(tokenString)
		if err != nil {
			//resource.RenderUnauthorizedError(w, req, "Unable to verify token string")
			return ctx
		}
		claims := c.(*TokenClaims)
	}

	return context.WithValue(ctx, TokenClaimsKey, c)
}
*/
