package sessions

import (
	"crypto/rsa"
	. "github.com/app/sessions/domain"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func generateJTI() string {
	// We will use mongodb's object id as JTI
	// we then will use this id to blacklist tokens,
	// along with `exp` and `iat` claims.
	// As far as collisions go, ObjectId is guaranteed unique
	// within a collection; and this case our collection is `sessions`
	return bson.NewObjectId().Hex()
}

// TokenAuthority implements ITokenAuthority
type TokenAuthority struct {
	Options *TokenAuthorityOptions
}

type TokenAuthorityOptions struct {
	PrivateSigningKeyByte []byte
	PublicSigningKeyByte  []byte
	PrivateSigningKey     *rsa.PrivateKey
	PublicSigningKey      *rsa.PublicKey
}

func NewTokenAuthority(options *TokenAuthorityOptions) *TokenAuthority {
	var err error
	options.PrivateSigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(options.PrivateSigningKeyByte)
	if err != nil {
		return nil
	}

	options.PublicSigningKey, err = jwt.ParseRSAPublicKeyFromPEM(options.PublicSigningKeyByte)
	if err != nil {
		return nil
	}

	ta := TokenAuthority{options}
	return &ta
}

func (ta *TokenAuthority) CreateNewSessionToken(claims ITokenClaims) (string, error) {

	c := claims.(*TokenClaims)

	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims = jwt.MapClaims{
		"userId": c.UserID,
		"exp":    time.Now().Add(time.Hour * 72).Format(time.RFC3339), // 3 days
		"iat":    time.Now().Format(time.RFC3339),
		"jti":    generateJTI(),
	}
	tokenString, err := token.SignedString(ta.Options.PrivateSigningKey)

	//DEBUG
	//fmt.Println(err)

	return tokenString, err
}

func (ta *TokenAuthority) VerifyTokenString(tokenString string) (IToken, ITokenClaims, error) {

	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ta.Options.PublicSigningKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	var claims TokenClaims
	tc := t.Claims.(jwt.MapClaims)
	token := NewToken(t)
	if tc.Valid() == nil {
		if tc["userId"] != nil {
			claims.UserID = tc["userId"].(string)
		}
		if tc["jti"] != nil {
			claims.JTI = tc["jti"].(string)
		}
		if tc["iat"] != nil {
			claims.IssuedAt, _ = time.Parse(time.RFC3339, tc["iat"].(string))
		}
		if tc["exp"] != nil {
			claims.ExpireAt, _ = time.Parse(time.RFC3339, tc["exp"].(string))
		}
	}

	return token, &claims, err
}
