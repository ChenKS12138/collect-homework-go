package admin

import (
	"log"
	"path/filepath"
	"time"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/email"
	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/template"
	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// serviceLogin servcie login
func serviceLogin(loginDto *LoginDto) (dataResponse *util.DataResponse,errResponse *util.ErrResponse){
	admin,err := database.Store.Admin.SelectByEmail(loginDto.Email);
	if admin == nil || (bcrypt.CompareHashAndPassword([]byte(admin.Password),[]byte(loginDto.Password)) != nil) {
		return nil,ErrAuthorization
	} 
	if err != nil {
		return nil,util.ErrRender(err)
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
	return util.NewDataResponse(claim.ToJwtClaim(time.Hour*4)),nil
}

// serviceInvitationCode service invitation code
func serviceInvitationCode(invitationCodeDto *InvitationCodeDto) (dataResponse *util.DataResponse,errResponse *util.ErrResponse){
	lastInvitationCode,err := database.Store.InvitationCode.SelectByEmail(invitationCodeDto.Email)

	if lastInvitationCode != nil && lastInvitationCode.CreateAt.Add(time.Minute).After(time.Now()) {
		return nil,ErrInvitationCodeFrequently
	}
	if err != nil {
		return nil,util.ErrRender(err)
	}

	code := util.RandString(6);

	mailText,err := template.Registry(code,invitationCodeDto.Email,time.Now());

	if err != nil {
		return nil,util.ErrRender(err)
	}
	if err :=email.SendMail(viper.GetString("EMAIL_USER"),"New Registry",mailText) ; err != nil {
		log.Println(err)
	}

	invitationCode := &model.InvitationCode{
		Email: invitationCodeDto.Email,
		Code: code,
	}
	database.Store.InvitationCode.Insert(invitationCode)
	return util.NewDataResponse(true),nil
}

// serviceRegister service register
func serviceRegister(registerDto *RegisterDto) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	invitationCode,err := database.Store.InvitationCode.SelectByEmail(registerDto.Email)

	if invitationCode == nil || invitationCode.Code != registerDto.InvitationCode {
		return nil,ErrInvitationCodeWrong
	}

	if err !=nil {
		return nil,util.ErrRender(err)
	}

	adminHasSameEmail,err := database.Store.Admin.SelectByEmail(registerDto.Email)

	if err != nil {
		return nil,util.ErrRender(err)
	}

	if adminHasSameEmail != nil {
		return nil,ErrEmailUsed
	}

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(registerDto.Password),10)

	admin := &model.Admin{
		Email: registerDto.Email,
		Password: string(hashedPassword),
		Name:registerDto.Name,
		IsSuperAdmin: false,
	}
	
	if err := database.Store.Admin.Insert(admin); err !=nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(true),nil
}

// serviceStatus servcie status
func serviceStatus(claim *auth.Claim)(dataResponse *util.DataResponse,errResponse *util.ErrResponse){
	fileCount,err := database.Store.Submission.SelectCountByAdminID(claim.ID)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	projects,err := database.Store.Project.SelectByAdminID(claim.ID)
	if err != nil {
		return nil,util.ErrRender(err)
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
			return nil,util.ErrRender(err)
		}
		totalSize+=size
	}

	return util.NewDataResponse(&struct{
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
	}),nil

}

// serviceSubToken servcie subToken
func serviceSubToken(subTokenDto *SubTokenDto,claim *auth.Claim)(dataResponse *util.DataResponse,errResponse *util.ErrResponse){
	claim.AuthCode = subTokenDto.AuthCode
	return util.NewDataResponse(claim.ToJwtClaim((time.Duration(subTokenDto.Expire)*time.Minute))),nil
}