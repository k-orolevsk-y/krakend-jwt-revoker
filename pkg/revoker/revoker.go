package revoker

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Revoker struct {
	blacklist *BlackList

	key     string
	keyType string
}

var (
	ErrorTokenIsRevoked = fmt.Errorf("token is revoked")
)

func NewRevoker(key, keyType string) *Revoker {
	return &Revoker{
		blacklist: NewBlackList(),

		key:     key,
		keyType: keyType,
	}
}

func (r *Revoker) Add(data map[string]string) error {
	r.blacklist.Add(data)
	return nil
}

func (r *Revoker) Middleware(req *http.Request) error {
	var authString string
	if r.keyType == "cookie" {
		cookie, err := req.Cookie(r.key)
		if err != nil {
			return err
		}
		authString = cookie.Value
	} else {
		authString = req.Header.Get(r.key)
	}
	authString = strings.ReplaceAll(authString, "Bearer ", "")

	claims, err := r.getTokenClaims(authString)
	if err != nil {
		return nil
	}

	if r.blacklist.Test(*claims) {
		return ErrorTokenIsRevoked
	}

	return nil
}

func (r *Revoker) getTokenClaims(token string) (*CustomClaims, error) {
	tkn, _, err := new(jwt.Parser).ParseUnverified(token, &CustomClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

type CustomClaims map[string]interface{}

func (claims *CustomClaims) Valid() error {
	return nil
}
