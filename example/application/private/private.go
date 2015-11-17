package private

import (
	"github.com/corneldamian/httpway"
	"github.com/corneldamian/httpwaymid"
	"github.com/julienschmidt/httprouter"

	"fmt"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	ctx := httpway.GetContext(r)

	fmt.Fprintf(w, "logged in: %s", ctx.Session().Username())

	ctx.Log().Info("profile")
}

func Logout(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	ctx := httpway.GetContext(r)

	sess := ctx.Session().(*httpwaymid.Session)
	sess.SetUsername("")
	sess.SetAuth(false)

	fmt.Fprint(w, "logged out")

	ctx.Log().Info("logout")
}
