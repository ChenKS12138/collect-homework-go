package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

// Claim claim
type Claim struct {
	IsSuperAdmin bool `json:"isSuperAdmin"`
	Name string `json:"name"`
	ID string `json:"id"`
	Email string `json:"email"`
}

// ToJwtClaim to jwt claim
func (c *Claim) ToJwtClaim() string {
	jwtClaim := jwt.MapClaims{
		"IsSuperAdmin":c.IsSuperAdmin,
		"Name":c.Name,
		"ID":c.ID,
		"Email":c.Email,
	}

	jwtauth.SetExpiry(jwtClaim,time.Now().Add(time.Hour*4))
	jwtauth.SetIssuedNow(jwtClaim)
	_,tokenString,_ := TokenAuth.Encode(jwtClaim)
	return tokenString;
}

// NewClaim from jwtclaim
func NewClaim(jwtClaim map[string]interface{}) (*Claim) {
	return &Claim{
		Email: jwtClaim["Email"].(string),
		IsSuperAdmin: jwtClaim["IsSuperAdmin"].(bool),
		ID: jwtClaim["ID"].(string),
		Name: jwtClaim["Name"].(string),
	}
}

// GenerateClaim generate claim
func GenerateClaim(r *http.Request) (*Claim,error) {
	_,jwtClaim,err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil,err
	}
	claim := NewClaim(jwtClaim)
	return claim,nil
}