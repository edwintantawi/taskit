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

func (j *jwtx) GenerateAccessToken(userID entity.UserID) (string, time.Time, error) {
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

func (j *jwtx) GenerateRefreshToken(userID entity.UserID) (string, time.Time, error) {
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

func (j *jwtx) VerifyAccessToken(rawToken string) (entity.UserID, error) {
	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token alg")
		}
		return []byte(j.accessTokenKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := entity.UserID(claims["user_id"].(string))
		return userID, nil
	} else {
		return "", errors.New("invalid token")
	}
}
