package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

// Code auth code
type Code = uint32;


// Create W; Read R; Update R+W; Delete X+W;
const (
	// CodeFileX file excuse
	CodeFileX Code=0b1 << 0
	// CodeFileW file read
	CodeFileW Code=0b1 << 1
	// CodeFileR file write
	CodeFileR Code=0b1 << 2

	// CodeProjectX file excuse
	CodeProjectX Code=0b1 << 3
	// CodeProjectW project excuse
	CodeProjectW Code=0b1 << 4
	// CodeProjectR project write
	CodeProjectR Code=0b1 << 5
	
	// CodeAdminX admin excuse
	CodeAdminX Code=0b1 << 6
	// CodeAdminW admin read
	CodeAdminW Code=0b1 << 7
	// CodeAdminR admin write
	CodeAdminR Code=0b1 << 8
)


// Claim claim
type Claim struct {
	IsSuperAdmin bool `json:"isSuperAdmin"`
	Name string `json:"name"`
	ID string `json:"id"`
	Email string `json:"email"`
	AuthCode Code `json:"authCode"`
}

// ToJwtClaim to jwt claim
func (c *Claim) ToJwtClaim(expire time.Duration) string {
	jwtClaim := jwt.MapClaims{
		"IsSuperAdmin":c.IsSuperAdmin,
		"Name":c.Name,
		"ID":c.ID,
		"Email":c.Email,
		"AuthCode":c.AuthCode,
	}

	jwtauth.SetExpiry(jwtClaim,time.Now().Add(expire))
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
		AuthCode: Code(jwtClaim["AuthCode"].(float64)),
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

// VerifyAuthCode verify auth code
func VerifyAuthCode(src Code,target Code)bool{
	return src & target == target;
}