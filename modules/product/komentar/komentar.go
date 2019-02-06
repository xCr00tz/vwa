package komentar

import (
	"fmt"
	"log"
	"net/http"

	//	"html/template"
	"vwa/helper/middleware"
	"vwa/util/database"
	"vwa/util/render"
	"vwa/util/session"

	"github.com/julienschmidt/httprouter"
)

type Self struct{}

func New() *Self {
	return &Self{}
}

type Komentar struct {
	IDKomentar  int    `json:"idkomentar"`
	IsiKomentar string `json:"isikomentar"`
	IDUser      int    `json:"iduser"`
	Username    string `json:"username"`
}

type Res struct {
	Data []Komentar `json:"data"`
}

var s = session.New()
var db, _ = database.Connect()

func (self *Self) SetRouter(r *httprouter.Router) {
	/* register all router */

	mw := middleware.New() //implement middleware

	r.GET("/fetchkomentar", mw.LoggingMiddleware(mw.CapturePanic(FetchKomentarHandler)))
	r.GET("/verifyuser", mw.LoggingMiddleware(mw.CapturePanic(VerifyUserHandler)))
	r.POST("/postkomentar", mw.LoggingMiddleware(mw.CapturePanic(SaveKomentarHandler)))
}

func FetchKomentarHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	res := Res{}
	res.Data = GetKomentar()
	render.JSONRender(w, res)

}

func VerifyUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	type Resp struct {
		Body string `json:"body"`
	}
	uid := s.GetSession(r, "id")

	if !s.IsLoggedIn(r) {
		html := `<div class="row">
                  <div class="col-lg-10">
                  <div class="alert alert-warning">Silahkan <strong>login</strong> untuk mengirim kometar</div>
                  </div>
                  </div>`
		resp := Resp{}
		resp.Body = html
		render.JSONRender(w, resp)
	} else {
		html := fmt.Sprintf(`<div class="card my-4">
				<h5 class="card-header"><i class="fa fa-comments-o"></i> <strong>Tinggalkan Komentar</strong></h5>
				<div class="card-body">
				<div id="kmsg" class="alert alert-danger" style="display:none"></div>
				<form id="formkomentar" action="#" method="post">
					<div class="form-group">
					<textarea class="form-control" name="isikomentar" rows="3"></textarea>
					<input type="hidden" name="uid" value="%s">
					</div>
				</form>
				<button type="submit" id="savekomentar" class="btn btn-primary"><i class="fa fa-comments"></i> <strong>Kirim Komentar</strong></button>
				</div>
			</div>
			<script>
			$("#savekomentar").on('click',function(){
				var data = $("#formkomentar").serialize()
				$.post("/postkomentar", data)
					.done(function(data){
						/* $("#kmsg").text(data[0].body)
						$("#kmsg").show()
						$("#kmsg").delay(2000).fadeOut(); */
						window.location.reload()
					})
			  })
			</script
			
			`, uid)

		resp := Resp{}
		resp.Body = html
		render.JSONRender(w, resp)
	}
}

func SaveKomentarHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	isiKomentar := r.FormValue("isikomentar")
	//filter := template.HTMLEscapeString(isiKomentar)
	uid := r.FormValue("uid")

	ok := SaveKomentar(uid, isiKomentar)
	if !ok {
		resp := struct {
			Body string `json:"body"`
		}{
			Body: "gagal mengirim komentar",
		}
		render.JSONRender(w, resp)
	}
}

func GetKomentar() []Komentar {

	const query = `
			select k.id_komentar, k.isi_komentar, k.id_user, u.username from komentar as k, users as u where k.id_user = u.id
			`
	rows, err := db.Query(query)
	defer rows.Close()

	data := []Komentar{} //function return data

	for rows.Next() {

		komentar := Komentar{} //kontainer for all komentar
		//var isikomentar string
		err = rows.Scan(&komentar.IDKomentar, &komentar.IsiKomentar, &komentar.IDUser, &komentar.Username)
		//komentar.IsiKomentar = template.HTMLEscapeString(isikomentar)
		if err != nil {
			log.Println(err.Error())
		}
		data = append(data, komentar)
		err = rows.Err()
		if err != nil {
			log.Println(err.Error())
		}
	}
	return data
}

func SaveKomentar(uid string, isikometar string) bool {

	const sqlQuery = `INSERT INTO komentar(isi_komentar, id_user) VALUES($1, $2)`
	_, err := db.Exec(sqlQuery, isikometar, uid)
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}

func GetUsername(uid string) string {

	var username string

	const sqlQuery = `select username from users where id=$1`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(err.Error())
	}
	defer stmt.Close()
	err = db.QueryRow(uid).Scan(&username)
	return username
}
