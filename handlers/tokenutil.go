package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/dgrijalva/jwt-go"
)

var InvalidJWTokenError = errors.New("Token was invalid or unreadable")

func ParseTokenFromString(env *Environment, tokenAsString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		strings.TrimSpace(tokenAsString),
		&JwtTokenClaim{},
		func(*jwt.Token) (interface{}, error) {
			return env.TokenSigningKey, nil
		})
}

func CreateTokenAsString(
	env *Environment,
	userId models.UserId,
	durationTilExpiration time.Duration,
) (string, error) {
	claims := JwtTokenClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(durationTilExpiration).Unix(),
			Issuer:    "CerealNotes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(env.TokenSigningKey)
}

func getUserIdFromJwtToken(env *Environment, request *http.Request) (models.UserId, error) {
	cookie, err := request.Cookie(cerealNotesCookieName)
	if err != nil {
		return 0, err
	}

	token, err := ParseTokenFromString(env, cookie.Value)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JwtTokenClaim); ok && token.Valid {
		return claims.UserId, nil
	}

	return 0, InvalidJWTokenError
}
