package routing

import (
	"github.com/corneldamian/httpway"
	"github.com/julienschmidt/httprouter"

	"fmt"
	"net/http"

	"github.com/corneldamian/gowebservice/example/application/private"
	"github.com/corneldamian/gowebservice/example/application/public"
)

func GetRoutesFactory() map[string]func(*httpway.Router) {

	routesFactory := map[string]func(*httpway.Router){
		"public":  publicRoutesFactory,
		"private": privateRoutesFactory,
	}

	return routesFactory
}

func MiddlewaresFactory(router *httpway.Router) (routes map[string]*httpway.Router) {
	routes = make(map[string]*httpway.Router)

	routes["public"] = router

	routes["private"] = router.Middleware(func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
		ctx := httpway.GetContext(r)

		if !ctx.Session().IsAuth() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not authenticated")
			return
		}

		ctx.Next(w, r, pr)
	})

	return
}

func publicRoutesFactory(router *httpway.Router) {
	router.GET("/", public.Index)
	router.GET("/login/:username", public.Login)
}

func privateRoutesFactory(router *httpway.Router) {
	router.GET("/profile", private.Profile)
	router.GET("/logout", private.Logout)
}
