package auth

import (
	"github.com/go-chi/jwtauth"
)

// TokenAuth token auth
var TokenAuth *jwtauth.JWTAuth;

const secret = "secret";

func init(){
	TokenAuth = jwtauth.New("HS256",[]byte(secret),nil);
}