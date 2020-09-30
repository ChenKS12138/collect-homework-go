package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)


const (
	contentTypeJSON string = "application/json"
	contentTypeFormData string = "multipart/form-data"
)

var client *http.Client

func init(){
	client = &http.Client{}
}

// GetRequest get request
func GetRequest(url string,header *http.Header,ri interface{},response interface{}) error{
	req,err := http.NewRequest("GET",url,nil)
	if header!=nil {
		req.Header = *header
	}
	res,err := client.Do(req)
	if err != nil {
		return err
	}
	responseString,err := parseResponse(res)

	if err != nil {
		return nil
	}
	return json.Unmarshal([]byte(responseString),response)
}

// PostRequest post request
func PostRequest(url string,header *http.Header,ri interface{},response interface{}) error {
	jsonString,err := newJSONString(ri)

	if err != nil {
		return err
	}

	req,err := http.NewRequest("POST",url,strings.NewReader(jsonString))
	if header != nil {
		req.Header = *header
	}
	req.Header.Set("Content-Type",contentTypeJSON)
	res,err := client.Do(req)
	responseString,err := parseResponse(res)

	if err !=nil {
		return err
	}
	return json.Unmarshal([]byte(responseString),response)
}


func newJSONString(i interface{}) (string,error) {
	bytes,err := json.Marshal(i)
	if err != nil {
		return "",err
	}
	return string(bytes),nil
}

func parseResponse(r *http.Response) (string,error) {
	bytes,err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "",err
	}
	return string(bytes),nil
}

// BasicResponse basic response
type BasicResponse struct {
	Data interface{} `json:"data"`
	Success bool `json:"success"`           // user-level status message
	StatusText string `json:"status"`
	ErrorText string `json:"error"`  // application-level error message, for debugging
}