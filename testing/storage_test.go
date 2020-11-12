package api_test

import (
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/ChenKS12138/collect-homework-go/api/app/storage"
	fileBytes "github.com/ChenKS12138/collect-homework-go/testing/bytes"
	"github.com/ChenKS12138/collect-homework-go/testing/service"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/chenhg5/collection"
)

// POST /storage/upload
func TestStorageUpload(t *testing.T){
	token,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal("Project Insert Abnormal")
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorageUploadWrongExtensions(t *testing.T){
	token,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_extensions_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"zip","rar"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	ok,err := service.StorageUpload(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileNameExtensions.ErrorText){
		t.Fatal(errors.New("Test Storage Upload Wrong Wrong Extensions Fail"))
	}
}

func TestStorageUploadWrongFileName(t *testing.T){
	token,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_filename_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	ok,err := service.StorageUpload(Ts.URL,util.RandString(6),projectID,"B11111111.doc",fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileNamePattern.ErrorText){
		t.Fatal(errors.New("Test Storage Upload Wrong Extensions Fail"))
	}
}

func TestStorageUploadWrongSecret(t *testing.T){
	token,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	fileName := "B11111111-陈陈陈-实验1.doc"
	secret := util.RandString(6)
	_,err = service.StorageUpload(Ts.URL,secret,projectID,fileName,fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	
	ok,err := service.StorageUpload(Ts.URL,secret+"_wrong_suffix",projectID,fileName,fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileSecret.ErrorText){
		t.Fatal("Test Storage Upload Fail (Expect ErrFileSecret)")
	}
	_,err = service.StorageUpload(Ts.URL,secret,projectID,fileName,fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
		t.Fatal("Test Storage Upload Fail (Expect Overwriting Success)")
	}
}

func TestStorageFileCount(t *testing.T){
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,"B11111112-陈陈-实验1.doc",fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	_,count,err := service.StorageFileCount(Ts.URL,commonUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal(errors.New("Test Stroage File Count Fail"))
	}

	_,count,err = service.StorageFileCount(Ts.URL,superUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal(errors.New("Test Storage File Count Fail (Super User)"))
	}

}

func TestStorageFileList(t *testing.T){
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken2,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	// common user
	_,filelist,err := service.StorageFileList(Ts.URL,commonUserToken,projectID)
	if filelist== nil || 
		!collection.Collect(filelist).Contains(projectNames[0]) ||
		!collection.Collect(filelist).Contains(projectNames[1]) || 
		collection.Collect(filelist).Count()!=2 {
		t.Fatal(errors.New("Test Storage File List Fail (Common User)"))
	}

	// super user
	_,filelist,err = service.StorageFileList(Ts.URL,superUserToken,projectID)
	if filelist== nil || 
	!collection.Collect(filelist).Contains(projectNames[0]) ||
	!collection.Collect(filelist).Contains(projectNames[1]) || 
	collection.Collect(filelist).Count()!=2 {
		t.Fatal(errors.New("Test Storage File List Fail (Super User)"))
	}

	// common user2
	ok,_,err := service.StorageFileList(Ts.URL,commonUserToken2,projectID)
	if ok || !strings.Contains(err.Error(),storage.ErrProjectPremissionDenied.ErrorText) {
		t.Fatal(errors.New("Test Storage File List Fail (Super User2)"))
	}
}

func TestStorageDownload(t *testing.T){
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken2,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	_,err = service.StorageDownload(Ts.URL,commonUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.StorageDownload(Ts.URL,superUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}

	ok,err := service.StorageDownload(Ts.URL,commonUserToken2,projectID)
	if ok || ! strings.Contains(err.Error(),"Not Bytes Stream") {
		t.Fatal("Storage Download Fail")
	}
}

func TestStorageProjectSize(t *testing.T){
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken1 ,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken2,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,commonUserToken1,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken1)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	// common user
	_,size,err := service.StorageProjectSize(Ts.URL,commonUserToken1,projectID)
	if err!=nil || size ==0 {
		log.Println(err,size);
		t.Fatal(errors.New("Test Storage Project Size Fail (Common User)"))
	}

	// super user
	_,size,err = service.StorageProjectSize(Ts.URL,superUserToken,projectID)
	if err != nil || size == 0 {
		log.Println(err,size)
		t.Fatal(errors.New("Test Storage Project Size Fail (Super User)"))
	}

	// common user2
	ok,_,err := service.StorageProjectSize(Ts.URL,commonUserToken2,projectID)
	if ok || !strings.Contains(err.Error(),storage.ErrProjectPremissionDenied.ErrorText) {
		t.Fatal(errors.New("Test Storage File List Fail (Super User2)"))
	}
}


func TestDownloadSelectively(t *testing.T){
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken2,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	_,err = service.StorageDownloadSelectively(Ts.URL,commonUserToken,projectID,"1")
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.StorageDownloadSelectively(Ts.URL,superUserToken,projectID,"A")
	if err != nil {
		t.Fatal(err)
	}

	ok,err := service.StorageDownloadSelectively(Ts.URL,commonUserToken2,projectID,"3")
	if ok || ! strings.Contains(err.Error(),"Not Bytes Stream") {
		t.Fatal("Storage Download Fail")
	}
}