package setup

import(
	"net/http"
	"vwa/helper/middleware"
	"github.com/julienschmidt/httprouter"
)

type Self struct{}

func New() *Self {
	return &Self{}
}
func (self *Self) SetRouter(r *httprouter.Router) {
	/* register all router */

	mw := middleware.New() //implement middleware

	r.GET("/setup", mw.LoggingMiddleware(mw.CapturePanic(SetupHandler)))

}

func SetupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
