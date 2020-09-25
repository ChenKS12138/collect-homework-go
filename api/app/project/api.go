package project

import (
	"collect-homework-go/auth"
	"collect-homework-go/database"
	"collect-homework-go/model"
	"collect-homework-go/util"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// Router router
func Router()(*chi.Mux,error){
	r := chi.NewRouter()

	// protected router
	r.Group(func (r chi.Router){
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator)
		
		r.Get("/own",own)
		r.Post("/insert",insert)
		r.Post("/update",update)
		// 不使用真正的delete
		// r.Post("/delete",delete)
	})

	// public router
	r.Group(func(r chi.Router){
		r.Get("/",list)
	})
	return r,nil
}

func own(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	projects,err := database.Store.Project.SelectByAdminID(claim.ID);
	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,projects)
}

// insert
func insert(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	insertDto := &InsertDto{}
	render.DecodeJSON(r.Body,insertDto)
	if err = insertDto.validate(); err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	project := &model.Project{
		AdminID: claim.ID,
		FileNameExample: insertDto.FileNameExample,
		FileNameExtensions: insertDto.FileNameExtensions,
		FileNamePattern: insertDto.FileNamePattern,
		Name: insertDto.Name,
	}
	
	err = database.Store.Project.Insert(project)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

// update
func update(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	updateDto := &UpdateDto{}
	render.DecodeJSON(r.Body,updateDto)
	if err = updateDto.validate(); err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	lastProject,err := database.Store.Project.SelectByID(updateDto.ID)
	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if lastProject.AdminID != claim.ID {
		render.Render(w,r,ErrProjectPermission)
		return
	}
	lastProject.Name = updateDto.Name
	lastProject.FileNamePattern = updateDto.FileNamePattern
	lastProject.FileNameExtensions = updateDto.FileNameExtensions
	lastProject.FileNameExample = updateDto.FileNameExample
	lastProject.Usable = updateDto.Usable

	if err = database.Store.Project.Update(lastProject);err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

// delete !discard
func delete(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	deleteDto := &DeleteDto{}
	render.DecodeJSON(r.Body,deleteDto)
	if err = deleteDto.validate();err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if err = database.Store.Project.Delete(deleteDto.ID,claim.ID); err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

// list
func list(w http.ResponseWriter, r *http.Request){
	projects,err := database.Store.Project.SelectAllUsable()
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(struct {
		Projects []model.ProjectWithAdminName `json:"projects"`;
	}{
		Projects: *projects,
	}))
}