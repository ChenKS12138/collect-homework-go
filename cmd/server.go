package cmd

import (
	"github.com/ChenKS12138/collect-homework-go/api"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string){
		LoadConfig()
		serve();
	},
}

func serve(){
	srv,_ := api.NewServer();
	srv.Start();
}

func init(){
	RootCmd.AddCommand(serveCmd)
}