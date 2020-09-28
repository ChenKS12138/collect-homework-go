package admin

import (
	"net/http"
	"time"

	"collect-homework-go/auth"
	"collect-homework-go/database"
	"collect-homework-go/email"
	"collect-homework-go/model"
	"collect-homework-go/template"
	"collect-homework-go/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

// Router router
func Router()(*chi.Mux ,error){
	r := chi.NewRouter();
	
	// protected router
	r.Group(func(r chi.Router){
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator)
		
	})

	// public router
	r.Group(func(r chi.Router){
		r.Post("/login",login)
		r.Post("/register",register)
		r.Post("/invitationCode",invitationCode)
	})
	return r,nil
}

// login
func login(w http.ResponseWriter,r *http.Request) {
	loginDto := &LoginDto{}
	render.DecodeJSON(r.Body,loginDto)
	if err := loginDto.validate(); err != nil {
		render.Render(w,r,util.ErrValidation(err))
		return
	} 
	
	admin,_ := database.Store.Admin.SelectByEmail(loginDto.Email);
	if admin == nil {
		render.Render(w,r,ErrAuthorization)
		return
	} 
	
	if bcrypt.CompareHashAndPassword([]byte(admin.Password),[]byte(loginDto.Password)) != nil {
		render.Render(w,r,ErrAuthorization)
		return
	}
		
	claim := &auth.Claim{
		IsSuperAdmin: admin.IsSuperAdmin,
		Email: admin.Email,
		ID: admin.ID,
		Name: admin.Name,
	}
	render.JSON(w,r,util.NewDataResponse(claim.ToJwtClaim()))
}

// invitation code
func invitationCode(w http.ResponseWriter,r *http.Request){
	invitationCodeDto := &InvitationCodeDto{}
	render.DecodeJSON(r.Body,invitationCodeDto)
	if err := invitationCodeDto.validate(); err != nil {
		render.Render(w,r,util.ErrValidation(err))
		return
	}

	if lastInvitationCode,err := database.Store.InvitationCode.SelectByEmail(invitationCodeDto.Email) ; err != nil || (lastInvitationCode!=nil && lastInvitationCode.CreateAt.Before(time.Now().Add(time.Minute))){
		if err!=nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			render.Render(w,r,ErrInvitationCodeFrequently)
		}
		return
	}

	code := util.RandString(6);

	mailText,err := template.Registry(code,invitationCodeDto.Email,time.Now());

	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if err :=email.SendMail( email.User,"New Registry",mailText,"Admin") ; err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	invitationCode := &model.InvitationCode{
		Email: invitationCodeDto.Email,
		Code: code,
	}
	database.Store.InvitationCode.Insert(invitationCode)
	render.JSON(w,r,util.NewDataResponse(true))
}

// register
func register(w http.ResponseWriter, r *http.Request){
	registerDto := &RegisterDto{}
	render.DecodeJSON(r.Body,registerDto)
	if err := registerDto.validate(); err !=nil {
		render.Render(w,r,util.ErrValidation(err))
		return
	}

	invationCode,err := database.Store.InvitationCode.SelectByEmail(registerDto.Email)
	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if invationCode == nil || invationCode.Code != registerDto.InvitationCode {
		render.Render(w,r,ErrInvitationCodeWrong)
		return
	}

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(registerDto.Password),10)

	admin := &model.Admin{
		Email: registerDto.Email,
		Password: string(hashedPassword),
		Name:registerDto.Name,
		IsSuperAdmin: false,
	}
	
	if err := database.Store.Admin.Insert(admin); err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}