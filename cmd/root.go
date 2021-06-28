package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/afocus/captcha"
	"github.com/coreos/etcd/pkg/fileutil"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tipFileName = "DON_NOT_TOUCH_THIS_DIR!!"

// ConfigFile config path
var ConfigFile string

// RootCmd root commmand
var RootCmd = &cobra.Command{
	Use:     "github.com/ChenKS12138/collect-homework-go",
	Short:   "Collect HomeWork Backend Build Time " + util.BuildTime,
	Version: util.Version,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config-file", "", "Env Config File Path")
}

// Execute execute
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// LoadConfig load config
func LoadConfig() {
	pwd, _ := os.Getwd()
	viper.AutomaticEnv()
	viper.SetConfigFile(filepath.Join(pwd, ConfigFile))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	secret := viper.GetString("JWT_SECRET") + util.Version
	auth.TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	storagePathPrefix := viper.GetString("STORAGE_PATH_PREFIX")
	tipFilePath := filepath.Join(storagePathPrefix, tipFileName)
	if !fileutil.Exist(tipFilePath) {
		ioutil.WriteFile(tipFilePath, []byte(tipFileName), 0664)
	}

	fontPath := viper.GetString("FONT_PATH")
	util.CaptchaCap = captcha.New()
	util.CaptchaCap.SetFont(fontPath)
	util.CaptchaCap.SetSize(256, 128)
	util.CaptchaCap.SetDisturbance(captcha.HIGH)
	util.CapachaSecret = secret
}
