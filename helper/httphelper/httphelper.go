package httphelper

import (
	"fmt"
	"net/http"

	"vwa/util"
)

func Redirect(w http.ResponseWriter, r *http.Request, location string, code int){
	redirect := fmt.Sprintf("%s%s", util.Fullurl,location)
	http.Redirect(w,r,redirect,code)
}