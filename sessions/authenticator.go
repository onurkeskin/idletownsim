package sessions

import (
	"github.com/app/users"
	//"github.com/app/server/domain"
	//. "github.com/app/sessions/domain"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

func NewAuthenticator(resource *Resource) *Authenticator {
	return &Authenticator{resource}
}

// SessionsAuthenticator implements IMiddleware
type Authenticator struct {
	resource *Resource
}

// Handler authenticates a session token in the Authorization header
func (auth *Authenticator) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	resource := auth.resource

	authHeaderString := req.Header.Get("Authorization")
	ctx := req.Context()

	if authHeaderString != "" {
		tokens := strings.Split(authHeaderString, " ")
		if len(tokens) != 2 || (len(tokens) > 0 && strings.ToUpper(tokens[0]) != "BEARER") {
			resource.RenderUnauthorizedError(w, req, "Invalid format, expected Authorization: Bearer {token}")
			return
		}

		tokenString := tokens[1]
		t, c, err := resource.TokenAuthority.VerifyTokenString(tokenString)
		if err != nil {
			resource.RenderUnauthorizedError(w, req, "Unable to verify token string")
			return
		}
		token := t.(*Token)
		claims := c.(*TokenClaims)

		if !token.Valid {
			resource.RenderUnauthorizedError(w, req, "Token is invalid")
			return
		}

		if time.Now().After(c.GetExpireAt()) {
			resource.RenderUnauthorizedError(w, req, "Token is expired")
			return
		}

		// Check that the token was not previously revoked
		// TODO: Possible optimization, use Redis
		revokedTokenRepo := resource.RevokedTokenRepository(req)
		if revokedTokenRepo.IsTokenRevoked(claims.JTI) {
			resource.RenderUnauthorizedError(w, req, "Token has been revoked")
			return
		}

		// retrieve user object and store it in current session request context
		// this `user` object will be used by the AccessController middleware
		userRepo := resource.UserRepository(req)
		user, err := userRepo.GetUserById(claims.UserID)
		if err != nil {
			// `user` = nil indicates that current authentication failed
			user = nil
		}

		// add token authority, claims and current user object to session request context
		//newContextWithTokenAuthority(ctx, req, token)
		//newContextWithAuthenticatedClaims(ctx, req, claims)
		ctx = context.WithValue(ctx, TokenAuthorityKey, token)
		ctx = context.WithValue(ctx, TokenClaimsKey, claims)
		ctx = users.NewContextWithUser(ctx, req, user)
	}
	/*else {
		fmt.Println("No Auth Data")
	}
	*/
	next(w, req.WithContext(ctx))
}

/*
func setTokenAuthorityCtx(ctx domain.IContext, r *http.Request, ta ITokenAuthority) {
	ctx.Set(r, TokenAuthorityKey, ta)
}
func setAuthenticatedClaimsCtx(ctx domain.IContext, r *http.Request, claim ITokenClaims) {
	ctx.Set(r, TokenClaimsKey, claim)
}
*/
