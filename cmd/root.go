package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigFile config path
var ConfigFile string

// RootCmd root commmand
var RootCmd = &cobra.Command{
	Use: "collect-homework-go",
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
}