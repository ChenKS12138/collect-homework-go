package service

import (
	"collect-homework-go/testing/request"
	"errors"
	"log"
)

// StorageUplaod storage upload
func StorageUplaod(baseURL string,secret string,projectID string,fileName string,fileBytes []byte) (ok bool,err error){
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