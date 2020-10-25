package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/spf13/viper"
)

// Server server
type Server struct {
	*http.Server
}

// NewServer new server
func NewServer() (*Server, error) {
	log.Println("Configuring Server...");
	apiRouter,_ := New();
	
	port := viper.GetString("port");

	srv := http.Server{
		Addr: ":"+port,
		Handler: apiRouter,
	}
	return &Server{&srv},nil;
}

// Start start
func (srv *Server) Start() {
	log.Println("Version: ",util.Version)
	log.Println("BuildTime: ",util.BuildTime)
	log.Println("Start Server...");
	go func(){
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err);
		}
	}()
	log.Printf("Listening on %s\n",srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit,os.Interrupt)
	sig := <- quit
	log.Println("shutting down server ... Reason:",sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}