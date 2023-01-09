package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID entity.UserID `json:"user_id"`
}

type JWTTokenConfig struct {
	Key string
	Exp int
}

type JWT struct {
	accessTokenKey      string
	refreshTokenKey     string
	accessTokenExpires  int
	refreshTokenExpires int
}

func NewJWT(accessToken, refreshToken JWTTokenConfig) JWT {
	return JWT{
		accessTokenKey:      accessToken.Key,
		refreshTokenKey:     refreshToken.Key,
		accessTokenExpires:  accessToken.Exp,
		refreshTokenExpires: refreshToken.Exp,
	}
}

func (j *JWT) GenerateAccessToken(userID entity.UserID) (string, time.Time, error) {
	expiresTime := time.Now().Add(time.Duration(j.accessTokenExpires) * time.Second)
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.accessTokenKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresTime, nil
}

func (j *JWT) GenerateRefreshToken(userID entity.UserID) (string, time.Time, error) {
	expiresTime := time.Now().Add(time.Duration(j.refreshTokenExpires) * time.Second)
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.refreshTokenKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresTime, nil
}

func (j *JWT) VerifyAccessToken(rawToken string) (entity.UserID, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token alg")
		}
		return []byte(j.accessTokenKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid access token")
	}

	userID := entity.UserID(claims["user_id"].(string))
	return userID, nil
}
