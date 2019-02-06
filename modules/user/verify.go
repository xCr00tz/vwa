package user

import (
	"fmt"
	"log"

	"net/http"
	"vwa/util"
	"vwa/util/database"
	"vwa/util/render"
	"vwa/util/session"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Body   string `json:"body"`
	Atas   string `json:"atas"`
	Status bool   `json:"status"`
}

func LoginVerify(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := session.New()
	if s.IsLoggedIn(r) {

		uid := s.GetSession(r, "id")

		data, err := GetProfile(uid)

		if err != nil {
			log.Println(err.Error())
		}

		res := Response{}

		profile := fmt.Sprintf(`
			<p><strong>Name :</strong> %s</p>
			<p><strong>Email :</strong> %s</p>
			<p><strong>Phone Number :</strong> %s</p>`,
			data.uname, data.email, data.msisdn)
		logout_atas := fmt.Sprintf(`<a class="nav-link" href="logout">Logout</a>`)
		res.Body = profile
		res.Atas = logout_atas
		res.Status = true
		render.JSONRender(w, res)

	} else {
		loginForm := fmt.Sprintf(`
		<div class="alert alert-danger" id="msg" style="display:none"></div>	
		<form id="loginform" method="post" action="#" accept-charset="utf-8">
			<fieldset>
				<div class="form-group">
					<input type="text" id="email"  name="email" value="" class="form-control" placeholder="Email" />
				</div>
				<div class="form-group">
					<input type="password" id="password" name="password" value="" class="form-control" placeholder="Password" />
				</div>
				
			</fieldset>
		</form> 
		<button id="btnlogin" class="btn btn-success btn-small btn-block"><i class="fa fa-sign-in"></i> <strong>Log in</strong></button>
		<script>
			function proses_login(){
				var loginurl = "%s/login"
				data = $('#loginform').serialize()
              	$.post(loginurl, data)
              	.done(function(res){
					if(res[0].success == false){
						$("#msg").text(res[0].body)
						$("#msg").show()
						$("#msg").delay(2000).fadeOut();
					}else{
						window.location.reload()
					}
            	}) 
			}
            $("#btnlogin").on('click',function(){
              proses_login();
          	})
			$("#email, #password").on("keypress", function (e) { 
			 var email 	= $("#email").val(),
			 	 pass 	= $("#password").val();
			 if(e.keyCode == 13){
			    if(email != "" && pass != ""){
			    	proses_login();
			    }
			  }
			}); 
		</script>`, util.Fullurl)
		logoutTop := fmt.Sprintf(`<a class="nav-link" href="#"></a>`)
		res := Response{}
		res.Body = loginForm
		res.Atas = logoutTop
		res.Status = false
		render.JSONRender(w, res)
	}
}

func GetProfile(uid string) (*UserData, error) {
	const (
		query = `SELECT username, email, phone_number FROM users where id=$1`
	)
	DB, _ := database.Connect()
	userdata := UserData{}

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(uid).Scan(&userdata.uname, &userdata.email, &userdata.msisdn)
	if err != nil {
		return nil, err
	}
	return &userdata, nil
}
