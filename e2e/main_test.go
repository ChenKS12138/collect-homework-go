package api_test

import (
	"collect-homework-go/api"
	"collect-homework-go/database/migrate"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

var (
	ts *httptest.Server
)

const (
	contentTypeJson string = "application/json"
)

func init(){

	pwd,_ := os.Getwd()
	viper.SetConfigFile(filepath.Join(pwd,"../.env"))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.AutomaticEnv()


	srv,_ := api.NewServer()
	ts = httptest.NewServer(srv.Handler)

	// clean dababase
	migrate.Init()
	migrate.Reset()
	migrate.Migrate([]string{})
}

// GET /
func TestWelcome(t *testing.T){
	response,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	bytes,err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(bytes),"Welcome!\nRequest From ") {
		t.Fatal(errors.New("Wrong Welcome Format"))
	}
	t.Log("Test Welcome Pass")
}

func getRequest(url string,ri interface{}) error{
	res,err := http.Get(url)
	if err != nil {
		return err
	}
	bytes,err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes,ri)
}

func postRequest(url string,ri interface{},response interface{}) error {
	jsonString,err := newJsonString(ri)

	if err != nil {
		return err
	}

	res,err := http.Post(url,contentTypeJson,strings.NewReader(jsonString))
	responseString,err := parseResponse(res)

	if err !=nil {
		return err
	}
	return json.Unmarshal([]byte(responseString),response)
}


func newJsonString(i interface{}) (string,error) {
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

type BasicResponse struct {
	Data interface{} `json:"data"`
	Success bool `json:"success"`           // user-level status message
	StatusText string `json:"status"`
	ErrorText string `json:"error"`  // application-level error message, for debugging
}