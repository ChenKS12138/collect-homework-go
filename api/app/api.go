package app

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/api/app/admin"
	"github.com/ChenKS12138/collect-homework-go/api/app/common"
	"github.com/ChenKS12138/collect-homework-go/api/app/project"
	"github.com/ChenKS12138/collect-homework-go/api/app/storage"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/go-chi/chi"
)

// Router router
func Router() (*chi.Mux, error) {

	r := chi.NewRouter()
	adminRouter, _ := admin.Router()
	projectRouter, _ := project.Router()
	storageRouter, _ := storage.Router()
	commonRouter, _ := common.Router()

	r.Get("/", welcome)
	r.Mount("/admin", adminRouter)
	r.Mount("/project", projectRouter)
	r.Mount("/storage", storageRouter)
	r.Mount("/common", commonRouter)

	return r, nil
}

func welcome(w http.ResponseWriter, r *http.Request) {
	text := "Welcome!\nRequest From " + r.RemoteAddr + "\nVersion: " + util.Version
	w.Write([]byte(text))
}
