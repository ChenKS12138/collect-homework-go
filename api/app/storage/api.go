package storage

import (
	"net/http"
	"strconv"

	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const maxFileSize = 100 * 1000 *1000
const textDot = "."

const (
	statusCreate string="新建"
	statusAlter string="修改"
)

// Router router
func Router()(*chi.Mux,error){
	r := chi.NewRouter()

	// protected router
	r.Group(func(c chi.Router){
		c.Use(jwtauth.Verifier(auth.TokenAuth))
		c.Use(jwtauth.Authenticator)
		
		// require auth.CodeFileR + auth.CodeFileX
		c.Get("/download",download)

		// require auth.CodeFileR
		c.Get("/fileList",fileList)

		// rquire auth.CodeProjectR
		c.Get("/projectSize",projectSize)
	})

	// public router
	r.Group(func(c chi.Router){
		c.Post("/upload",upload)
		c.Get("/fileCount",fileCount)
	})
	return r,nil
}

// upload
func upload(w http.ResponseWriter,r *http.Request){
	// 入参检验以及文件大小检验
	file,fileHeader,err := r.FormFile("file")
	ip := r.RemoteAddr
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else {
		uploadDto := &UploadDto{
			File: file,
			FileHeader: fileHeader,
			ProjectID: r.FormValue("projectId"),
			Secret: r.FormValue("secret"),
		}
		if err := uploadDto.validate(); err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceUpload(uploadDto,ip)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}

func download(w http.ResponseWriter,r *http.Request){
	// 入参检查
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeFileR+auth.CodeFileX) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		values:= r.URL.Query()
		downloadDto := &DownloadDto{
			ID: values.Get("id"),
		}
		if err := downloadDto.validate(); err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,filename,err := serviceDownload(downloadDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				w.Header().Set("Content-Length",strconv.FormatInt(int64(len(*data)),10))
				w.Header().Set("Content-Disposition",`attachment;filename="`+filename+`.zip"`)
				render.Data(w,r,*data)
			}
		}
	}
}

func fileCount(w http.ResponseWriter,r *http.Request){
	values := r.URL.Query()
	fileCountDto := &FileCountDto{
		ID: values.Get("id"),
	}
	if err := fileCountDto.validate();err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else {
		data,err := serviceFileCount(fileCountDto)
		if err != nil {
			render.Render(w,r,err)
		} else {
			render.Render(w,r,data)
		}
	}
}

func fileList(w http.ResponseWriter,r *http.Request){
	// 入参检查
	claim,err := auth.GenerateClaim(r)
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeFileR) {
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		values := r.URL.Query()
		fileListDto := &FileListDto{
			ID: values.Get("id"),
		}
		if err := fileListDto.validate();err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceFileList(fileListDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}

func projectSize(w http.ResponseWriter,r *http.Request){
	claim,err := auth.GenerateClaim(r);
	if err != nil {
		render.Render(w,r,util.ErrRender(err))
	} else if !auth.VerifyAuthCode(claim.AuthCode,auth.CodeProjectR){
		render.Render(w,r,util.ErrUnauthorized)
	} else {
		values := r.URL.Query()
		projectSizeDto := &ProjectSizeDto{
			ID: values.Get("id"),
		}
		if err := projectSizeDto.validate();err != nil {
			render.Render(w,r,util.ErrRender(err))
		} else {
			data,err := serviceProjectSize(projectSizeDto,claim)
			if err != nil {
				render.Render(w,r,err)
			} else {
				render.Render(w,r,data)
			}
		}
	}
}