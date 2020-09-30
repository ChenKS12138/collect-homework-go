package api_test

import (
	"collect-homework-go/api"
	"collect-homework-go/auth"
	"collect-homework-go/database/migrate"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var (
	ts *httptest.Server
	superAdmin struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
)


func init(){

	pwd,_ := os.Getwd()
	viper.SetConfigFile(filepath.Join(pwd,"../.env"))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.AutomaticEnv()

	superAdmin.Email = viper.GetString("SUPER_USER_EMAIL")
	superAdmin.Password = viper.GetString("SUPER_USER_PASSWORD")
	auth.TokenAuth =  jwtauth.New("HS256",[]byte(viper.GetString("JWT_SECRET")),nil)

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
	// t.Log("Test Welcome Pass")
}
