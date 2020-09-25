package app

import (
	"net/http"

	"collect-homework-go/api/app/admin"
	"collect-homework-go/api/app/project"
	"collect-homework-go/api/app/storage"

	"github.com/go-chi/chi"
)

// Router router
func Router() (*chi.Mux,error) {

	r:=chi.NewRouter()
	adminRouter,_ := admin.Router()
	projectRouter,_ := project.Router();
	storageRouter,_ := storage.Router();
	
	r.Get("/",welcome)
	r.Mount("/admin",adminRouter)
	r.Mount("/project",projectRouter)
	r.Mount("/storage",storageRouter)

	return r,nil
}


func welcome(w http.ResponseWriter,r *http.Request){
	text := "Welcome!\nRequest From "+r.RemoteAddr
	w.Write([]byte(text))
}