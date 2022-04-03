package cmd

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ChenKS12138/collect-homework-go/api"
	"github.com/ChenKS12138/collect-homework-go/resource"
	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/go-chi/chi"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		LoadConfig()
		serve()
	},
}

// Server server
type Server struct {
	*http.Server
}

// NewServer new server
func NewServer() (*Server, error) {
	log.Println("Configuring Server...")
	apiRouter, _ := api.New()
	r := chi.NewRouter()

	publicFS, err := fs.Sub(resource.Public, "public")
	if err != nil {
		panic(err)
	}
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
	r.Mount("/api", apiRouter)
	r.Mount("/", http.FileServer(http.FS(publicFS)))

	port := viper.GetString("port")

	srv := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	return &Server{&srv}, nil
}

// Start start
func (srv *Server) Start() {
	log.Println("Version: ", util.Version)
	log.Println("BuildTime: ", util.BuildTime)
	log.Println("Start Server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("shutting down server ... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}

func serve() {
	srv, _ := NewServer()
	srv.Start()
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
