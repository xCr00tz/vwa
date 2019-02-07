package profile

import (
	"fmt"
	"log"
	"net/http"

	"crypto/md5"
	"encoding/hex"
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

func (self *Self) SetRouter(r *httprouter.Router) {

	mw := middleware.New()

	r.GET("/verify_user", mw.LoggingMiddleware(mw.CapturePanic(Verify_User)))
	r.GET("/user", mw.LoggingMiddleware(mw.CapturePanic(UserHandler)))
	r.POST("/user", mw.LoggingMiddleware(mw.CapturePanic(GetUserHandler)))
	r.GET("/profile", mw.LoggingMiddleware(mw.CapturePanic(ProfileHandler)))
	r.POST("/profile", mw.LoggingMiddleware(mw.CapturePanic(UpdateProfileHandler)))
	r.POST("/password", mw.LoggingMiddleware(mw.CapturePanic(UpdatePasswordHandler)))

}

var DB, _ = database.Connect()

type UserData struct {
	UserID   string `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	MSISDN   string `json:"msisdn"`
}

type Jsonresp struct {
	Success string    `json:"success"`
	Data    *UserData `json:"data"`
	Message string    `json:"message"`
}

func UserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	data := make(map[string]interface{})
	nama := r.URL.Query()["user"][0]
	data["title"] = "User Profile"
	data["nama_user"] = nama

	render.HTMLRender(w, r, "template.user", data)

}

func GetUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		uid := r.FormValue("uid")
		respdata, err := GetUserData(uid)
		if err != nil {
			resp := Jsonresp{}
			resp.Success = "0"
			resp.Data = respdata
			resp.Message = err.Error()
			render.JSONRender(w, resp)
		} else {
			resp := Jsonresp{}
			resp.Success = "1"
			resp.Data = respdata
			resp.Message = ""
			render.JSONRender(w, resp)
		}
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	sess := session.New()
	data := make(map[string]interface{})

	if sess.IsLoggedIn(r) {
		uid := sess.GetSession(r, "id")
		userdata, err := GetProfile(uid)

		if err != nil {
			log.Println(err.Error())
		}

		data["title"] = "Profile"
		data["uid"] = userdata.UserID
		data["email"] = userdata.Email
		data["name"] = userdata.UserName
		data["msisdn"] = userdata.MSISDN
	} else {
		data["title"] = "Profile"
	}

	render.HTMLRender(w, r, "template.profile", data)

}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := session.New()
	resp := Jsonresp{}

	if sess.IsLoggedIn(r) {
		if r.Method == "POST" {
			uid := sess.GetSession(r, "id")
			name := r.FormValue("name")
			email := r.FormValue("email")
			msisdn := r.FormValue("msisdn")
			ok := updateProfile(uid, name, email, msisdn)
			if !ok {
				resp.Success = "0"
				resp.Message = "Gagal menperbaharui data"
			} else {
				resp.Success = "1"
				resp.Message = "Data berhasil diperbaharui"
			}
		}
	} else {
		resp.Message = "0"
		resp.Message = "Login untuk dapat memperbaharui data"
	}
	render.JSONRender(w, resp)
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := session.New()
	resp := Jsonresp{}

	if sess.IsLoggedIn(r) {
		if r.Method == "POST" {
			uid := r.FormValue("uid")
			password_lama := r.FormValue("password_lama")
			ok := updatePassword(uid, password_lama)
			if !ok {
				resp.Success = "0"
				resp.Message = "Gagal Mengganti Password"
			} else {
				resp.Success = "1"
				resp.Message = "Password Berhasil Diganti"
			}
		}
	} else {
		resp.Message = "0"
		resp.Message = "Login untuk dapat memperbaharui data"
	}
	render.JSONRender(w, resp)
}

func GetUserData(uid string) (*UserData, error) {

	query := fmt.Sprintf("SELECT username, email, phone_number FROM users where id=%s", uid)
	userdata := UserData{}
	stmt := DB.QueryRow(query)

	err := stmt.Scan(&userdata.UserName, &userdata.Email, &userdata.MSISDN)
	if err != nil {
		return nil, err
	}
	return &userdata, nil
}

func GetProfile(uid string) (*UserData, error) {
	const (
		query = `SELECT username, email, phone_number FROM users where id=$1`
	)
	userdata := UserData{}

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(uid).Scan(&userdata.UserName, &userdata.Email, &userdata.MSISDN)
	if err != nil {
		return nil, err
	}
	return &userdata, nil
}

func updateProfile(uid string, name string, email string, phone_number string) bool {
	const (
		query = `UPDATE users SET username=$1, email=$2, phone_number=$3 where id = $4`
	)
	_, err := DB.Exec(query, name, email, phone_number, uid)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func Md5Sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func updatePassword(uid string, password_baru string) bool {
	const (
		query = `UPDATE users SET password=$1 where id = $2`
	)
	_, err := DB.Exec(query, Md5Sum(password_baru), uid)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

type Resp struct {
	Body string `json:"body"`
}

func Verify_User(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := session.New()
	if s.IsLoggedIn(r) == false {

		html := `<div class="alert alert-warning">Silahkan <strong>login</strong> untuk melihat halaman ini</div>`
		resp := Resp{}
		resp.Body = html
		render.JSONRender(w, resp)

	}
}
