package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd root commmand
var RootCmd = &cobra.Command{
	Use: "collect-homework-go",
}

// Execute execute
func Execute(){
	viper.AutomaticEnv(); // read in environment variables that match
	if err := RootCmd.Execute(); err !=nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
