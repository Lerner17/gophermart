package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lerner17/gophermart/internal/config"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateAccessToken(user *models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	cfg := config.Instance

	return generateToken(user, expirationTime, []byte(cfg.JWTSecretKey))
}

func GenerateTokensAndSetCookies(user *models.User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)

	return nil
}

func generateToken(user *models.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		ID:       user.ID,
		Username: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *models.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Login
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

var ErrUnauthorized = &er.HTTPError{
	Code: 401,
	Msg:  "unauthorized",
}

func JWTErrorChecker(err error, c echo.Context) error {
	return ErrUnauthorized
}

func GetUserIDFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cfg := config.Instance
		return []byte(cfg.JWTSecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*models.JwtCustomClaims); ok && token.Valid {
		return claims.ID, nil
	} else {
		return 0, fmt.Errorf("cannot parse jwt")
	}

}
