package user

import (

	"log"
	"net/http"
	"strconv"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"vwa/util/render"
	"vwa/helper/httphelper"
	"vwa/util/session"
	"vwa/util/database"
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

	r.GET("/login", mw.LoggingMiddleware(mw.CapturePanic(LoginHandler)))
	r.GET("/verify", mw.LoggingMiddleware(mw.CapturePanic(LoginVerify)))
	r.POST("/login", mw.LoggingMiddleware(mw.CapturePanic(LoginHandler)))
	r.GET("/logout", mw.LoggingMiddleware(Logout))
}

func LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	type Resp struct{
		Success bool `json:"success"`
		Body string `json:"body"`
	}
	if r.Method == "POST" {

		if !validateForm(w,r,ps) {
			res := Resp{}
			res.Success = false
			res.Body = "Empty email or Password"
			render.JSONRender(w, res)
		}else{
			if loginAction(w, r, ps) {
				res := Resp{}
				res.Success = true
				res.Body = ""
				render.JSONRender(w, res)
			} else {
				res := Resp{}
				res.Success = false
				res.Body = "Incorrect email or password"
				render.JSONRender(w, res)
				log.Println("Login Failed")
			}
		}
	}
}

func loginAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) bool {

	/* handler for login action */
	uname := r.FormValue("email")
	pass := Md5Sum(r.FormValue("password"))

	uData := checkUserQuery(uname, pass) //handle user data from db
	if uData.cnt == 1 {
		s := session.New()

		/* save user data to session */
		sessionData := make(map[string]string)
		sessionData["id"] = strconv.Itoa(uData.id)
		sessionData["uname"] = uData.uname
		sessionData["email"] = uData.email
		sessionData["msisdn"] = uData.msisdn

		s.SetSession(w, r, sessionData)
		//util.SetCookie(w, "Uid", strconv.Itoa(uData.id)) //save user_id to cookie
		log.Println("Login Success")

		return true
	} else {
		return false
	}
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s := session.New()
	s.DeleteSession(w, r)
	httphelper.Redirect(w, r, "index", 302)
}

func validateForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params)bool{
	uname := r.FormValue("email")
	pass := r.FormValue("password")
	if uname == "" || pass == ""{
		return false
	}
	return true
}

/* type to handle user data that return form query */
type UserData struct {
	id    int
	uname string
	email string
	msisdn string
	cnt int
}

var db *sql.DB

func checkUserQuery(username, pass string) *UserData {
	/* this function will check rows num which return from query */
	db, err := database.Connect()
	if err != nil {
		log.Println(err.Error())
	}

	var uData = UserData{} //inisialize empty userdata

	const (
		sql = `SELECT id, username, email, phone_number, COUNT(*) as cnt
						FROM users 
						WHERE email=$1 
						AND password=$2
						group by id,username,email,phone_number`)

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(username, pass).Scan(&uData.id, &uData.uname, &uData.email, &uData.msisdn, &uData.cnt)
	return &uData

}

func Md5Sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
