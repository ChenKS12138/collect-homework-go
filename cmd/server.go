package cmd

import (
	"collect-homework-go/api"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string){
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