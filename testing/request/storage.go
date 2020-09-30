package request

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

// StorageUpload storage upload
func StorageUpload(url string,projectID string,secret string,fileName string, fileBytes []byte)(*struct {
	BasicResponse
	Data bool `json:"data"`
},error){
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("projectId",projectID)
	bodyWriter.WriteField("secret",secret)
	fileWriter,err := bodyWriter.CreateFormFile("file",fileName)
	if err != nil {
		return nil,err
	}
	_,err = io.Copy(fileWriter,bytes.NewReader(fileBytes))
	if err != nil {
		return nil,err
	}
	bodyWriter.Close()
	r,err := http.Post(url,bodyWriter.FormDataContentType(),bodyBuf)
	responseString,err := parseResponse(r)
	if err != nil {
		return nil,err
	}
	response := &struct {
		BasicResponse
		Data bool `json:"data"`
	}{}
	err = json.Unmarshal([]byte(responseString),response)
	return response,err
}

// StorageFileCount storage file count
func StorageFileCount(url string,token string,projectID string) (* struct {
	BasicResponse
	Data struct {
		Count int `json:"count"`
	} `json:"data"`
},error) {
	response := &struct {
		BasicResponse
		Data struct {
			Count int `json:"count"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	err := GetRequest(url,header,map[string]string{
		"id":projectID,
	},response)
	return response,err
}

// StorageFileList storage file list
func StorageFileList(url string,token string,projectID string)(*struct {
	BasicResponse
	Data struct {
		Files []string `json:"files"`
	} `json:"data"`
},error) {
	response := &struct {
		BasicResponse
		Data struct {
			Files []string `json:"files"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	err := GetRequest(url,header,map[string]string{
		"id":projectID,
	},response)
	return response,err
}