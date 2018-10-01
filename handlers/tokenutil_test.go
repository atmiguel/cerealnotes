package handlers_test

import (
	"fmt"
	"testing"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/test_util"
)

func TestToken(t *testing.T) {
	env := &handlers.Environment{nil, []byte("TheWorld")}

	var num models.UserId = 32
	bob, err := handlers.CreateTokenAsString(env, num, 1)
	if err != nil {
		panic(err)
	}

	token, err := handlers.ParseTokenFromString(env, bob)
	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(*handlers.JwtTokenClaim); ok && token.Valid {
		test_util.Equals(t, int64(32), int64(claims.UserId))

		fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println("Token claims could not be read")
		t.FailNow()
	}
}
