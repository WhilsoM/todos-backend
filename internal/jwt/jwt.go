package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
	}
}

type Claims struct {
	UserID int `json:"user_id"`

	jwtlib.RegisteredClaims
}

func (j *JWTManager) GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTManager) ParseToken(tokenString string) (int, error) {
	token, err := jwtlib.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtlib.Token) (any, error) {
			return j.secret, nil
		},
		jwtlib.WithValidMethods([]string{
			jwtlib.SigningMethodHS256.Alg(),
		}),
	)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}
