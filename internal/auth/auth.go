package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

type authKey string

const AuthKey authKey = "authToken"

func NewAuth(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*Auth, error) {
	if publicKey == nil || privateKey == nil {
		return &Auth{}, errors.New("private/public key cannot be nil")
	}
	return &Auth{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (a *Auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}
	return tokenStr, nil
}

func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims

	tkn, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) { return a.publicKey, nil })
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}

	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}

	return c, nil
}
