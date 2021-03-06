package service

import (
	"errors"
	"log"

	"github.com/ChenKS12138/collect-homework-go/testing/request"
)

// StorageUpload storage upload
func StorageUpload(baseURL string,secret string,projectID string,fileName string,fileBytes []byte) (ok bool,err error){
	response,err := request.StorageUpload(baseURL+"/storage/upload",projectID,secret,fileName,fileBytes)
	if err != nil {
		return false,err
	}
	if !response.Success {
		log.Println(response)
		return false,errors.New(response.ErrorText)
	}
	return true,nil
}

// StorageFileCount storage file count
func StorageFileCount(baseURL string,token string,projectID string) (ok bool,count int,err error){
	response,err := request.StorageFileCount(baseURL+"/storage/fileCount",token,projectID)
	if err != nil {
		return false,0,err
	}
	if !response.Success {
		log.Println(response)
		return false,0,errors.New(response.ErrorText)
	}
	return true,response.Data.Count,nil
}

// StorageFileList storage file list
func StorageFileList(baseURL string,token string,projectID string)(ok bool,filelist []string,err error){
	response,err := request.StorageFileList(baseURL+"/storage/fileList",token,projectID)
	if err != nil {
		return false,nil,err
	}
	if !response.Success {
		return false,nil,errors.New(response.ErrorText)
	}
	return true,response.Data.Files,nil
}

// StorageDownload storage download
func StorageDownload(baseURL string,token string,projectID string) (ok bool,err error){
	return request.StorageDownload(baseURL+"/storage/download",token,projectID)
}

// StorageProjectSize storage project size
func StorageProjectSize(baseURL string,token string,projectID string)(ok bool,size int64,err error){
	response,err := request.StorageProjectSize(baseURL+"/storage/projectSize",token,projectID)
	if err != nil {
		return false,0,err
	}
	if !response.Success {
		log.Println(response)
		return false,0,errors.New(response.ErrorText)
	}
	return true,response.Data.Size,nil
}

// StorageDownloadSelectively storage download selectively
func StorageDownloadSelectively(baseURL string,token string,projectID string,code string) (ok bool,err error){
	return request.StorageDownloadSelectively(baseURL+"/storage/downloadSelectively",token,projectID,code)
}

// StorageRawFile storage raw file
func StorageRawFile(baseURL string,token string,projectID string,filename string) (ok bool,err error){
	return request.StorageRawFile(baseURL+"/storage/rawFile",token,projectID,filename)
}