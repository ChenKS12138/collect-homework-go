package admin

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// Router router
func Router()(*chi.Mux ,error){
	r := chi.NewRouter()
	
	// protected router
	r.Group(func(c chi.Router){
		c.Use(jwtauth.Verifier(auth.TokenAuth))
		c.Use(jwtauth.Authenticator)
		
		// rquire auth.CodeAdminR + auth.CodeProjectR + auth.CodeFileR
		c.Get("/status",status)
		// require [x| Authoriry x higher than SubTokenDto.AuthCode] 
		c.Post("/subToken",subToken)
	})

	// public router
	r.Group(func(c chi.Router){
		c.Post("/login",login)
		c.Post("/register",register)
		c.Post("/invitationCode",invitationCode)
	})
	return r,nil
}

// login
func login(w http.ResponseWriter,r *http.Request) {
	loginDto := &LoginDto{}
	render.DecodeJSON(r.Body,loginDto)
	if err := loginDto.validate(); err != nil {
		render.Render(w,r,util.ErrValidation(err))
	} else {
		data,err := serviceLogin(loginDto)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}
	
}

// invitation code
func invitationCode(w http.ResponseWriter,r *http.Request){
	invitationCodeDto := &InvitationCodeDto{}
	render.DecodeJSON(r.Body,invitationCodeDto)
	if err := invitationCodeDto.validate(); err != nil {
		render.Render(w,r,util.ErrValidation(err))
	} else {
		data,err := serviceInvitationCode(invitationCodeDto)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}

	
}

// register
func register(w http.ResponseWriter, r *http.Request){
	registerDto := &RegisterDto{}
	render.DecodeJSON(r.Body,registerDto)
	if err := registerDto.validate(); err !=nil {
		render.Render(w,r,util.ErrValidation(err))
	} else {
		data,err := serviceRegister(registerDto)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}
	
}

// status
func status(w http.ResponseWriter, r *http.Request) {
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		return
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeAdminR+auth.CodeProjectR+auth.CodeFileR) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		data,err := serviceStatus(claim)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}
}

//subToken
func subToken(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else {
		subTokenDto := &SubTokenDto{};
		render.DecodeJSON(r.Body,subTokenDto);
	
		if err:= subTokenDto.validate(); err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			if !auth.VerifyAuthCode(claim.AuthCode,subTokenDto.AuthCode) {
				render.Render(w,r,ErrInsufficientAuthority)
				return
			}
			data,err := serviceSubToken(subTokenDto,claim)
			if err != nil {
				render.Render(w,r,data)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}