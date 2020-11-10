package admin

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/email"
	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/template"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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
		c.Post("/sign",sign)
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
		return
	} 
	
	admin,err := database.Store.Admin.SelectByEmail(loginDto.Email);
	if admin == nil || (bcrypt.CompareHashAndPassword([]byte(admin.Password),[]byte(loginDto.Password)) != nil) {
		render.Render(w,r,ErrAuthorization)
		return
	} 
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	//full auth
	authCode := auth.CodeFileX +
							auth.CodeFileW +
							auth.CodeFileR + 
							auth.CodeProjectX +
							auth.CodeProjectW +
							auth.CodeProjectR +
							auth.CodeAdminX +
							auth.CodeAdminW +
							auth.CodeAdminR
		
	claim := &auth.Claim{
		IsSuperAdmin: admin.IsSuperAdmin,
		Email: admin.Email,
		ID: admin.ID,
		Name: admin.Name,
		AuthCode: authCode,
	}

	render.JSON(w,r,util.NewDataResponse(claim.ToJwtClaim(time.Hour*4)))
}

// invitation code
func invitationCode(w http.ResponseWriter,r *http.Request){
	invitationCodeDto := &InvitationCodeDto{}
	render.DecodeJSON(r.Body,invitationCodeDto)
	if err := invitationCodeDto.validate(); err != nil {
		render.Render(w,r,util.ErrValidation(err))
		return
	}

	lastInvitationCode,err := database.Store.InvitationCode.SelectByEmail(invitationCodeDto.Email)

	if lastInvitationCode != nil && lastInvitationCode.CreateAt.Add(time.Minute).After(time.Now()) {
		render.Render(w,r,ErrInvitationCodeFrequently)
		return
	}
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	code := util.RandString(6);

	mailText,err := template.Registry(code,invitationCodeDto.Email,time.Now());

	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if err :=email.SendMail(viper.GetString("EMAIL_USER"),"New Registry",mailText) ; err != nil {
		log.Println(err)
		// render.Render(w,r,util.ErrRender(err))
		// return
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

	invitationCode,err := database.Store.InvitationCode.SelectByEmail(registerDto.Email)

	if invitationCode == nil || invitationCode.Code != registerDto.InvitationCode {
		render.Render(w,r,ErrInvitationCodeWrong)
		return
	}

	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	adminHasSameEmail,err := database.Store.Admin.SelectByEmail(registerDto.Email)

	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if adminHasSameEmail != nil {
		render.Render(w,r,ErrEmailUsed)
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

// status
func status(w http.ResponseWriter, r *http.Request) {
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeAdminR+auth.CodeProjectR+auth.CodeFileR) {
		render.Render(w,r,util.ErrUnauthorized)
		return
	}

	fileCount,err := database.Store.Submission.SelectCountByAdminID(claim.ID)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	projects,err := database.Store.Project.SelectByAdminID(claim.ID)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	projectCount := len(*projects)

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	fileutil.TouchDirAll(filepath.Join(storagePathPrefix))

	totalSize := int64(0)
	for _,project := range(*projects) {
		dirPath := filepath.Join(storagePathPrefix,project.ID) 
		fileutil.TouchDirAll(dirPath)
		size,err := util.DirSizeB(dirPath)
		if err != nil {
			render.Render(w,r,util.ErrRender(err))
			return
		}
		totalSize+=size
	}

	render.JSON(w,r,util.NewDataResponse(&struct{
		FileCount int `json:"fileCount"`
		ProjectCount int `json:"projectCount"`
		TotalSize int64 `json:"totalSize"`
		Username string `json:"username"`
		Email string `json:"email"`
	} {
		FileCount: fileCount,
		ProjectCount: projectCount,
		TotalSize: totalSize,
		Username: claim.Name,
		Email: claim.Email,
	}))
}

//sign
func sign(w http.ResponseWriter,r *http.Request){
	// claim,err := auth.GenerateClaim(r)
	// if err != nil {
	// 	render.Render(w,r,util.ErrRender(err))
	// 	return
	// }
	// values:=r.URL.Query()
	// expire,_ := strconv.Atoi(values.Get("expire"));
	// signDto := &SignDto{
	// 	Expire: expire,
	// }
	// if err:= signDto.validate(); err != nil {
	// 	render.Render(w,r,util.ErrRender(err))
	// 	return
	// }
	
	// authCode := auth.CodeProjectFile + 
	// 	auth.CodeProjectRead +
	// 	auth.CodeProjectExcuse
	// claim.AuthCode = authCode

	// render.JSON(w,r,util.NewDataResponse(claim.ToJwtClaim((time.Minute*10))));
}