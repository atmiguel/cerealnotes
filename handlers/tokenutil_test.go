package handlers_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/models"
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
	fmt.Println(bob)
	if claims, ok := token.Claims.(*handlers.JwtTokenClaim); ok && token.Valid {
		if claims.UserId != 32 {
			fmt.Println("error in token")
			t.FailNow()
		}
		fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println("Token claims could not be read")
		t.FailNow()
	}
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
