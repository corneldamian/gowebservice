package public

import (
	"github.com/corneldamian/httpway"
	"github.com/corneldamian/httpwaymid"

	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := httpway.GetContext(r)

	sess := ctx.Session().(*httpwaymid.Session)
	sess.SetUsername(ctx.ParamByName("username"))
	sess.SetAuth(true)

	fmt.Fprint(w, "logged in")

	ctx.Log().Info("login")
}

func Index(w http.ResponseWriter, r *http.Request) {
	ctx := httpway.GetContext(r)

	if ctx.Session().IsAuth() {
		fmt.Fprintf(w, "index data for: %s", ctx.Session().Username())
	} else {
		fmt.Fprint(w, "index data")
	}

	ctx.Log().Info("index")
}
