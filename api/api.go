package api

import (
	"net/http"
	"time"

	"github.com/ChenKS12138/collect-homework-go/api/app"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// New new
func New() (*chi.Mux,error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15*time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers	
	}).Handler)

	_,err := database.DBConn();
	if err !=nil {
		panic(err);
	}

	appRouter,_ := app.Router()
	r.Mount("/",appRouter)
	r.MethodNotAllowed(routeUnable)
	r.NotFound(routeUnable)

	return r,nil
}

func routeUnable(w http.ResponseWriter,r *http.Request){
	render.Render(w,r,util.ErrNotFound)
}