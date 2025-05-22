package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignedTokenParams struct {
	AccountId string
	Nounce    string
}

type JwtPayload struct {
	AccountId uuid.UUID
	ExpiresAt int64
	jwt.RegisteredClaims
}
