package security

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type tokenclaims struct {
	jwt.RegisteredClaims
	Data map[string]any `json:"data"`
}

type JWTTokenConfig struct {
	Key string
	Exp int
}

type jwtx struct {
	accessTokenKey      string
	refreshTokenKey     string
	accessTokenExpires  int
	refreshTokenExpires int
}

func NewJWT(accessToken, refreshToken JWTTokenConfig) *jwtx {
	return &jwtx{
		accessTokenKey:      accessToken.Key,
		refreshTokenKey:     refreshToken.Key,
		accessTokenExpires:  accessToken.Exp,
		refreshTokenExpires: refreshToken.Exp,
	}
}

func (j *jwtx) GenerateAccessToken(payload map[string]any) (string, time.Time, error) {
	expiresTime := time.Now().Add(time.Duration(j.accessTokenExpires) * time.Second)
	claims := tokenclaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
		payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.accessTokenKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresTime, nil
}

func (j *jwtx) GenerateRefreshToken(payload map[string]any) (string, time.Time, error) {
	expiresTime := time.Now().Add(time.Duration(j.refreshTokenExpires) * time.Second)
	claims := tokenclaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
		payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.refreshTokenKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresTime, nil
}
