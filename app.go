package main

import (
	"fmt"
	"net/http"
	"vwa/helper/middleware"
	"vwa/modules/product/komentar"
	product "vwa/modules/product/main"
	"vwa/modules/user"
	"vwa/modules/user/profile"
	"vwa/util"
	"vwa/util/render"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	data["title"] = "Home"
	render.HTMLRender(w, r, "template.index", data)
}

func main() {
	mw := middleware.New()
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.GET("/", mw.LoggingMiddleware(indexHandler))
	router.GET("/index", mw.LoggingMiddleware(indexHandler))

	user := user.New()
	user.SetRouter(router)

	komentar := komentar.New()
	komentar.SetRouter(router)

	product := product.New()
	product.SetRouter(router)

	profile := profile.New()
	profile.SetRouter(router)

	s := http.Server{
		Addr:    ":" + util.Cfg.Webport,
		Handler: router,
	}

	fmt.Printf("Server running at port %s\n", s.Addr)
	fmt.Printf("Open this url %s on your browser to access VWA", util.Fullurl)
	s.ListenAndServe()
}
