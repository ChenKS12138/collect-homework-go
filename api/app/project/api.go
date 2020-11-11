package project

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/util"

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
		
		// require auth.CodeProjectR
		c.Get("/own",own)
		
		// require auth.CodeProjectW
		c.Post("/insert",insert)
		// require auth.CodeProjectW + auth.CodeProjectR
		c.Post("/update",update)
		
		// 不使用真正的delete
		// require auth.CodeProjectW + auth.CodeProjectX
		c.Post("/delete",delete)

		// require auth.CodeProjectW + auth.CodeProjectX
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
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectR){
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		data,err := serviceOwn(claim)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}
	
}

// insert
func insert(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectW) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		insertDto := &InsertDto{}
		render.DecodeJSON(r.Body,insertDto)
		if err = insertDto.validate(); err !=nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceInsert(insertDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}

// update
func update(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectW+auth.CodeProjectR) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		updateDto := &UpdateDto{}
		render.DecodeJSON(r.Body,updateDto)
		if err = updateDto.validate(); err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceUpdate(updateDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}


// delete
func delete(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else	if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectW+auth.CodeProjectX) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		deleteDto := &DeleteDto{}
		render.DecodeJSON(r.Body,deleteDto)
		if err = deleteDto.validate();err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceDelete(deleteDto,claim);
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}

// restore
func restore(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectW+auth.CodeProjectX) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		restoreDto := &RestoreDto{}
		render.DecodeJSON(r.Body,restoreDto)
		if err = restoreDto.validate();err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceRestore(restoreDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}

// list
func list(w http.ResponseWriter, r *http.Request){
	data,err := serviceList()
	if err != nil {
		render.Render(w,r,err)
	} else {
		render.Render(w,r,data)
	}
}