package storage

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/email"
	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/template"
	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/chenhg5/collection"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func serviceUpload(uploadDto *UploadDto,ip string) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	if uploadDto.FileHeader.Size > maxFileSize {
		return nil,ErrFileSize
	}

	// project存在检验
	lastProject,err := database.Store.Project.SelectAdminEmailByID(uploadDto.ProjectID)
	if lastProject == nil {
		return nil,ErrProjectNotExist
	}
	if err !=nil  {
		return nil,util.ErrRender(err)
	}

	// 判断是否为第一次上传
	lastSubmission,err := database.Store.Submission.SelectByProjectIDAndName(lastProject.ID,uploadDto.FileHeader.Filename)
	if err != nil{
		return nil,util.ErrRender(err)
	}
	if lastSubmission != nil && (bcrypt.CompareHashAndPassword([]byte(lastSubmission.Secret),[]byte(uploadDto.Secret)) !=nil ){
		return nil,ErrFileSecret
	}

	// 文件名/扩展名 正则检验
	fileNamePrefix := uploadDto.FileHeader.Filename
	if len(lastProject.FileNameExample)!=0 {
		fileNameExtensionCollections := collection.Collect(lastProject.FileNameExtensions)
		if fileNameExtensionCollections.Count() > 0 {
			dotIndex := strings.LastIndex(fileNamePrefix,textDot)
			if !fileNameExtensionCollections.Contains(fileNamePrefix[(dotIndex+1):]) {
				return nil,ErrFileNameExtensions
			}
			fileNamePrefix = fileNamePrefix[:(dotIndex)]
		}
	}
	if len(lastProject.FileNamePattern)!=0 {
		ok,err := regexp.Match(lastProject.FileNamePattern,[]byte(fileNamePrefix))
		if err!=nil {
			return nil,util.ErrRender(err)
		}
		if !ok {
			return nil,ErrFileNamePattern
		}
	}
	

	// 文件写入&记录存储&邮件发送
	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	if !fileutil.Exist(storagePathPrefix) {
		fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	}
	dirPath := filepath.Join(storagePathPrefix,lastProject.ID)
	if !fileutil.Exist(dirPath) {
		fileutil.TouchDirAll(dirPath)
	}
	filePath := filepath.Join(dirPath,uploadDto.FileHeader.Filename)
	fileBytes,err := ioutil.ReadAll(uploadDto.File)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	if err = ioutil.WriteFile(filePath,fileBytes,0664);err !=nil {
		return nil,util.ErrRender(err)
	}
	secret,err := bcrypt.GenerateFromPassword([]byte(uploadDto.Secret),10)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	m:= md5.New()
	m.Write(fileBytes);
	md5Str := hex.EncodeToString(m.Sum(nil))
	abspath,err:=filepath.Abs(filePath)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	submission := &model.Submission{
		FileName: uploadDto.FileHeader.Filename,
		IP: ip,
		ProjectID: uploadDto.ProjectID,
		FilePath: abspath,
		Secret: string(secret),
		MD5: md5Str,
	}
	if err= database.Store.Submission.Insert(submission);err !=nil {
		return nil,util.ErrRender(err)
	}
	if lastProject.SendEmail {
		statusText :=statusCreate
		if lastSubmission != nil {
			statusText = statusAlter
		}
		mailText,err := template.Submission(lastProject.Name,statusText,uploadDto.FileHeader.Filename,time.Now(),ip,md5Str)
		if err!=nil {
			return nil,util.ErrRender(err)
		}
		err = email.SendMail(lastProject.AdminEmail,"New Submission",mailText)
		if err!=nil {
			return nil,util.ErrRender(err)
		}
	}
	return util.NewDataResponse(true),nil
}

func serviceDownload(downloadDto *DownloadDto,claim *auth.Claim)(bytes *[]byte,filename string,errResponse *util.ErrResponse){
	// project 存在性检查
	project,err := database.Store.Project.SelectByID(downloadDto.ID)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	if project == nil || (!claim.IsSuperAdmin && claim.ID != project.AdminID ) {
		return nil,"",ErrDownloadForbidden
	}

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	tmpPathPrefix := viper.GetString("TEMP_PATH_PREFIX")
	if !fileutil.Exist(storagePathPrefix) {
		fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	}
	if !fileutil.Exist(tmpPathPrefix) {
		fileutil.TouchDirAll(filepath.Join(tmpPathPrefix))
	}
	dirPath := filepath.Join(storagePathPrefix,project.ID)
	if !fileutil.Exist(dirPath) {
		fileutil.TouchDirAll(dirPath)
	}
	zipFilePath := filepath.Join(tmpPathPrefix,project.Name+"-"+strconv.Itoa(int(time.Now().Unix()))+".zip")
	err = util.Zip(dirPath,zipFilePath)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	zipBytes,err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	return &zipBytes,project.Name,nil
}

func serviceFileCount(fileCountDto *FileCountDto) (dataResponse *util.DataResponse,errResponse *util.ErrResponse){
	count,err := database.Store.Submission.SelectCountByProjectID(fileCountDto.ID)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(&struct{
		Count int `json:"count"`
	}{
		Count: count,
	}),nil
}

func serviceFileList(fileListDto *FileListDto,claim *auth.Claim)  (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	if !claim.IsSuperAdmin {
		project,err := database.Store.Project.SelectByAdminIDAndID(claim.ID,fileListDto.ID)
		if err != nil {
			return nil,util.ErrRender(err)
		}
		if project == nil {
			return nil,ErrProjectPremissionDenied
		}
	}

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	if !fileutil.Exist(filepath.Join(storagePathPrefix)) {
		fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	}
	dirPath := filepath.Join(storagePathPrefix,fileListDto.ID)
	if !fileutil.Exist(dirPath) {
		fileutil.TouchDirAll(dirPath)
	}

	files,err := ioutil.ReadDir(dirPath)

	if err != nil {
		return nil,util.ErrRender(err)
	}

	filelist := []string{}

	for _,file := range(files) {
		if !file.IsDir() {
			filelist = append(filelist,file.Name())
		}
	}
	return util.NewDataResponse(&struct{
		Files []string `json:"files"`
	}{
		Files: filelist,
	}),nil
}

func serviceProjectSize(projectSizeDto *ProjectSizeDto,claim *auth.Claim)  (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	project,err := database.Store.Project.SelectByID(projectSizeDto.ID)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	if !claim.IsSuperAdmin && project.AdminID != claim.ID {
		return nil,ErrProjectPremissionDenied
	}

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	dirPath := filepath.Join(storagePathPrefix,project.ID)
	if !fileutil.Exist(dirPath) {
		fileutil.TouchDirAll(dirPath)
	}
	size,err := util.DirSizeB(dirPath)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(&struct{
		Size int64 `json:"size"`
	} {
		Size:size,
	}),nil
}

func serviceDownloadSelectively(downloadSelectivelyDto *DownloadSelectivelyDto,claim *auth.Claim)(bytes *[]byte,filename string,errResponse *util.ErrResponse){
	// project 存在性检查
	project,err := database.Store.Project.SelectByID(downloadSelectivelyDto.ID)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	if project == nil || (!claim.IsSuperAdmin && claim.ID != project.AdminID ) {
		return nil,"",ErrDownloadForbidden
	}

	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	tmpPathPrefix := viper.GetString("TEMP_PATH_PREFIX")
	if !fileutil.Exist(storagePathPrefix) {
		fileutil.TouchDirAll(filepath.Join(storagePathPrefix))
	}
	if !fileutil.Exist(tmpPathPrefix) {
		fileutil.TouchDirAll(filepath.Join(tmpPathPrefix))
	}
	// sourceDir := filepath.Join(storagePathPrefix,project.ID)
	tmpDir := filepath.Join(tmpPathPrefix,project.ID+"_"+downloadSelectivelyDto.Code)
	if !fileutil.Exist(tmpDir) {
		fileutil.TouchDirAll(filepath.Join(tmpDir))
	}
	
	downloadCode := util.ParseDownloadCode(downloadSelectivelyDto.Code)
	filePathMap,err := database.Store.Submission.SelectFilePathMap(downloadSelectivelyDto.ID)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	files,err := database.Store.Submission.SelectAllFile(downloadSelectivelyDto.ID)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	codeLen := len(*downloadCode)
	for fileIndex,file := range(*files) {
		if fileIndex < codeLen && (*downloadCode)[fileIndex] == 1 {
			sourceFilePath := filepath.Join((*filePathMap)[file.FileName])
			tmpFilePath := filepath.Join(tmpDir,file.FileName)
			if len(sourceFilePath) >0 && fileutil.Exist(sourceFilePath) {
				source, _ := os.Open(sourceFilePath)
        defer source.Close()
				destination, _ := os.Create(tmpFilePath)
        defer destination.Close()
				io.Copy(destination,source)
			}
		}
	}
	zipFilePath := filepath.Join(tmpPathPrefix,project.Name+"_"+downloadSelectivelyDto.Code+"_"+strconv.Itoa(int(time.Now().Unix()))+".zip")
	err = util.Zip(tmpDir,zipFilePath)
	if err != nil {
		return nil,"",util.ErrRender(err)
	}
	zipBytes,err := ioutil.ReadFile(zipFilePath)
	return &zipBytes,project.Name,nil
}