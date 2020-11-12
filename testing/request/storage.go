package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
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

// StorageDownload storage download
func StorageDownload(url string,token string,projectID string) (bool,error) {
	req,_ := http.NewRequest("GET",url,nil)
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	req.Header = *header
	q := req.URL.Query()
	q.Add("id",projectID)
	req.URL.RawQuery = q.Encode()
	res,err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false,err
	}
	if strings.Contains(res.Header.Get("Content-Type"),contentTypeJSON) {
		responseByte,_ := ioutil.ReadAll(res.Body)
		log.Println(string(responseByte))
		return false,errors.New("Not Bytes Stream")
	}
	return true,nil
}

// StorageProjectSize project size 
func StorageProjectSize(url string,token string,projectID string)(*struct {
	BasicResponse
	Data struct {
		Size int64 `json:"size"`
	} `json:"data"`
},error){
	response := &struct {
		BasicResponse
		Data struct {
			Size int64 `json:"size"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	err := GetRequest(url,header,map[string]string {
		"id":projectID,
	},response)
	return response,err
}

// StorageDownloadSelectively storage download selectively
func StorageDownloadSelectively(url string,token string,projectID string,code string) (bool,error) {
	req,_ := http.NewRequest("GET",url,nil)
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	req.Header = *header
	q := req.URL.Query()
	q.Add("id",projectID)
	q.Add("code",code)
	req.URL.RawQuery = q.Encode()
	res,err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false,err
	}
	if strings.Contains(res.Header.Get("Content-Type"),contentTypeJSON) {
		responseByte,_ := ioutil.ReadAll(res.Body)
		log.Println(string(responseByte))
		return false,errors.New("Not Bytes Stream")
	}
	return true,nil
}