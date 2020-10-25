package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile config path
var ConfigFile string

// RootCmd root commmand
var RootCmd = &cobra.Command{
	Use: "github.com/ChenKS12138/collect-homework-go",
	Short: "Collect HomeWork Backend Build Time "+util.BuildTime,
	Version: util.Version,
}

func init(){
	RootCmd.PersistentFlags().StringVar(&ConfigFile,"config-file","","Env Config File Path")
}

// Execute execute
func Execute(){
	if err := RootCmd.Execute(); err !=nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// LoadConfig load config
func LoadConfig(){
	pwd, _ := os.Getwd()
	viper.AutomaticEnv()
	viper.SetConfigFile(filepath.Join(pwd,ConfigFile))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	auth.TokenAuth =  jwtauth.New("HS256",[]byte(viper.GetString("JWT_SECRET")),nil)
}