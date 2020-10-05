package storage

import (
	"collect-homework-go/auth"
	"collect-homework-go/database"
	"collect-homework-go/email"
	"collect-homework-go/model"
	"collect-homework-go/template"
	"collect-homework-go/util"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chenhg5/collection"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const maxFileSize = 50 * 1000 *1000
const textDot = "."

const (
	statusCreate string="新建"
	statusAlter string="修改"
)

// Router router
func Router()(*chi.Mux,error){
	r := chi.NewRouter()

	// protected router
	r.Group(func(c chi.Router){
		c.Use(jwtauth.Verifier(auth.TokenAuth))
		c.Use(jwtauth.Authenticator)
		c.Get("/download",download)
		c.Get("/fileList",fileList)
	})

	// public router
	r.Group(func(c chi.Router){
		c.Post("/upload",upload)
		c.Get("/fileCount",fileCount)
	})
	return r,nil
}

// upload
func upload(w http.ResponseWriter,r *http.Request){
	// 入参检验以及文件大小检验
	file,fileHeader,err := r.FormFile("file")
	ip := r.RemoteAddr
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	uploadDto := &UploadDto{
		File: file,
		FileHeader: fileHeader,
		ProjectID: r.FormValue("projectId"),
		Secret: r.FormValue("secret"),
	}
	if err := uploadDto.validate(); err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if uploadDto.FileHeader.Size > maxFileSize {
		render.Render(w,r,ErrFileSize)
		return
	}

	// project存在检验
	lastProject,err := database.Store.Project.SelectAdminEmailByID(uploadDto.ProjectID)
	if lastProject == nil {
		render.Render(w,r,ErrProjectNotExist)
		return
	}
	if err !=nil  {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	// 判断是否为第一次上传
	lastSubmission,err := database.Store.Submission.SelectByProjectIDAndName(lastProject.ID,uploadDto.FileHeader.Filename)
	if err != nil{
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if lastSubmission != nil && (bcrypt.CompareHashAndPassword([]byte(lastSubmission.Secret),[]byte(uploadDto.Secret)) !=nil ){
		render.Render(w,r,ErrFileSecret)
		return
	}

	// 文件名/扩展名 正则检验
	fileNamePrefix := uploadDto.FileHeader.Filename
	fileNameExtensionCollections := collection.Collect(lastProject.FileNameExtensions)
	if fileNameExtensionCollections.Count() > 0 {
		dotIndex := strings.LastIndex(fileNamePrefix,textDot)
		if !fileNameExtensionCollections.Contains(fileNamePrefix[(dotIndex+1):]) {
			render.Render(w,r,ErrFileNameExtensions)
			return
		}
		fileNamePrefix = fileNamePrefix[:(dotIndex)]
	}
	ok,err := regexp.Match(lastProject.FileNamePattern,[]byte(fileNamePrefix))
	if err!=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if !ok {
		render.Render(w,r,ErrFileNamePattern)
		return
	}

	// 文件写入&记录存储&邮件发送
	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	dirPath := filepath.Join(storagePathPrefix,lastProject.ID)
	fileutil.TouchDirAll(dirPath)
	filePath := filepath.Join(dirPath,uploadDto.FileHeader.Filename)
	fileBytes,err := ioutil.ReadAll(uploadDto.File)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if err = ioutil.WriteFile(filePath,fileBytes,0664);err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	secret,err := bcrypt.GenerateFromPassword([]byte(uploadDto.Secret),10)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	submission := &model.Submission{
		FileName: uploadDto.FileHeader.Filename,
		IP: ip,
		ProjectID: uploadDto.ProjectID,
		FilePath: filePath,
		Secret: string(secret),
	}
	if lastSubmission == nil {
		if err= database.Store.Submission.Insert(submission);err !=nil {
			render.Render(w,r,util.ErrRender(err))
			return
		}
	}
	statusText :=statusCreate
	if lastSubmission != nil {
		statusText = statusAlter
	}
	mailText,err := template.Submission(lastProject.Name,statusText,uploadDto.FileHeader.Filename,time.Now(),ip)
	if err!=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	err = email.SendMail(lastProject.AdminEmail,"New Submission",mailText)
	if err!=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

func download(w http.ResponseWriter,r *http.Request){
	// 入参检查
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	values:= r.URL.Query()
	downloadDto := &DownloadDto{
		ID: values.Get("id"),
	}
	if err := downloadDto.validate(); err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	// project 存在性检查
	project,err := database.Store.Project.SelectByID(downloadDto.ID)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if project == nil || (!claim.IsSuperAdmin && claim.ID != project.AdminID ) {
		render.Render(w,r,ErrDownloadForbidden)
		return
	}

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	dirPath := filepath.Join(storagePathPrefix,project.ID)
	fileutil.TouchDirAll(dirPath)
	zipFilePath := filepath.Join(storagePathPrefix,project.Name+"-"+string(time.Now().Unix())+".zip")
	err = util.Zip(dirPath,zipFilePath)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	zipBytes,err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.Data(w,r,zipBytes)
}

func fileCount(w http.ResponseWriter,r *http.Request){
	values := r.URL.Query()
	fileCountDto := &FileCountDto{
		ID: values.Get("id"),
	}
	if err := fileCountDto.validate();err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	count,err := database.Store.Submission.SelectCountByProjectID(fileCountDto.ID)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(&struct{
		Count int `json:"count"`
	}{
		Count: count,
	}))
}

func fileList(w http.ResponseWriter,r *http.Request){
	// 入参检查
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	values := r.URL.Query()
	fileListDto := &FileListDto{
		ID: values.Get("id"),
	}
	if err := fileListDto.validate();err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	 if !claim.IsSuperAdmin {
		 project,err := database.Store.Project.SelectByAdminIDAndID(claim.ID,fileListDto.ID)
		 if err != nil {
			 render.Render(w,r,util.ErrRender(err))
			 return
		 }
		 if project == nil {
			 render.Render(w,r,ErrProjectPremissionDenied)
			 return
		 }
	 }

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	dirPath := filepath.Join(storagePathPrefix,fileListDto.ID)
	fileutil.TouchDirAll(dirPath)

	files,err := ioutil.ReadDir(dirPath)

	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	filelist := []string{}

	for _,file := range(files) {
		if !file.IsDir() {
			filelist = append(filelist,file.Name())
		}
	}
	render.JSON(w,r,util.NewDataResponse(&struct{
		Files []string `json:"files"`
	}{
		Files: filelist,
	}))
}