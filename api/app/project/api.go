package project

import (
	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/util"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// Router router
func Router()(*chi.Mux,error){
	r := chi.NewRouter()

	// protected router
	r.Group(func (c chi.Router){
		c.Use(jwtauth.Verifier(auth.TokenAuth))
		c.Use(jwtauth.Authenticator)
		
		c.Get("/own",own)
		c.Post("/insert",insert)
		c.Post("/update",update)
		// 不使用真正的delete
		c.Post("/delete",delete)
		c.Post("/restore",restore)
	})

	// public router
	r.Group(func(c chi.Router){
		c.Get("/",list)
	})
	return r,nil
}

// own
func own(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	projects := &[]model.ProjectWithAdminName{}

	// 超级管理员允许查看所有
	if claim.IsSuperAdmin {
		projects,err = database.Store.Project.SelectAllWithName()
	} else {
		projects,err = database.Store.Project.SelectByAdminID(claim.ID);
	}

	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(struct {
		Projects []model.ProjectWithAdminName `json:"projects"`;
	}{
		Projects: *projects,
	}))
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

	if lastProject == nil {
		render.Render(w,r,ErrProjectNotFound)
		return
	}

	if err !=nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}

	if !claim.IsSuperAdmin && lastProject.AdminID != claim.ID {
		render.Render(w,r,ErrProjectPermission)
		return
	}
	lastProject.FileNamePattern = updateDto.FileNamePattern
	lastProject.FileNameExtensions = updateDto.FileNameExtensions
	lastProject.FileNameExample = updateDto.FileNameExample
	lastProject.Usable = updateDto.Usable
	lastProject.SendEmail = updateDto.SendEmail
	lastProject.Visible = updateDto.Visible

	if err = database.Store.Project.Update(lastProject);err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}


// delete
func delete(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	deleteDto := &DeleteDto{}
	render.DecodeJSON(r.Body,deleteDto)
	if err = deleteDto.validate();err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	lastProject,err := database.Store.Project.SelectByID(deleteDto.ID)

	if lastProject == nil {
		render.Render(w,r,ErrProjectNotFound)
	}
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if !claim.IsSuperAdmin && lastProject.AdminID != claim.ID {
		render.Render(w,r,ErrProjectPermission)
		return
	}
	lastProject.Usable = false
	if err = database.Store.Project.Update(lastProject);err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

// restore
func restore(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	restoreDto := &RestoreDto{}
	render.DecodeJSON(r.Body,restoreDto)
	if err = restoreDto.validate();err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	lastProject,err := database.Store.Project.SelectByID(restoreDto.ID)

	if lastProject == nil {
		render.Render(w,r,ErrProjectNotFound)
	}
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	if !claim.IsSuperAdmin {
		render.Render(w,r,ErrProjectPermission)
		return
	}
	lastProject.Usable = true
	if err = database.Store.Project.Update(lastProject);err != nil {
		render.Render(w,r,util.ErrRender(err))
		return
	}
	render.JSON(w,r,util.NewDataResponse(true))
}

// delete !discard
// func delete(w http.ResponseWriter,r *http.Request){
// 	claim,err := auth.GenerateClaim(r)
// 	if err != nil {
// 		render.Render(w,r,util.ErrRender(err))
// 		return
// 	}

// 	deleteDto := &DeleteDto{}
// 	render.DecodeJSON(r.Body,deleteDto)
// 	if err = deleteDto.validate();err !=nil {
// 		render.Render(w,r,util.ErrRender(err))
// 		return
// 	}

// 	if err = database.Store.Project.Delete(deleteDto.ID,claim.ID); err != nil {
// 		render.Render(w,r,util.ErrRender(err))
// 		return
// 	}
// 	render.JSON(w,r,util.NewDataResponse(true))
// }

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