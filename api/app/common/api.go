package common

import (
	"image/png"
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/afocus/captcha"
	"github.com/go-chi/chi"
)

func Router() (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Group(func(c chi.Router) {
		c.Get("/generateCaptcha", generateCaptcha)
	})
	return r, nil
}

// generateCaptcha
func generateCaptcha(w http.ResponseWriter, r *http.Request) {
	img, str := util.CaptchaCap.Create(6, captcha.ALL)
	secret := util.GenerateCaptchaSecret()
	enc, err := util.Encrypt(secret, &str)
	if err != nil {
		panic(err)
	}
	w.Header().Add("X-Captcha", *enc)
	w.Header().Add("Content-Type", "image/png")
	png.Encode(w, img)
}
