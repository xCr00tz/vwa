package product

import (
	"log"
	"net/http"

	//	"html/template"
	"vwa/helper/middleware"
	"vwa/util/database"
	"vwa/util/render"

	"github.com/julienschmidt/httprouter"
)

type Self struct{}

func New() *Self {
	return &Self{}
}

var db, _ = database.Connect()

type Product struct {
	IDProduct        int    `json:"idproduct"`
	FotoProduct      string `json:"fotoproduct"`
	JudulProduct     string `json:"judulproduct"`
	DeskripsiProduct string `json:"deskripsiproduct"`
	Terjual          int    `json:"terjual"`
	Disukai          int    `json:"disukai"`
	Harga            int    `json:"harga"`
}

type Jsonresp struct {
	Success string    `json:"success"`
	Data    []Product `json:"data"`
	Message string    `json:"message"`
}

func (self *Self) SetRouter(r *httprouter.Router) {
	/* register all router */

	mw := middleware.New() //implement middleware

	r.GET("/product", mw.LoggingMiddleware(FetchProductPage))
	r.GET("/cari_product", mw.LoggingMiddleware(FetchListProduct))
}
func FetchProductPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	data["title"] = "List Product"
	render.HTMLRender(w, r, "template.product", data)
}

func FetchListProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	urutkan := r.FormValue("urutkan")
	berdasar := r.FormValue("berdasar")
	filter := r.FormValue("filter")
	pencarian := r.FormValue("pencarian")
	dari := r.FormValue("dari")
	hingga := r.FormValue("hingga")

	respdata, err := GetProductFilter(urutkan, berdasar, filter, pencarian, dari, hingga)
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

func GetProductFilter(urutkan, berdasar, filter, pencarian, dari, hingga string) ([]Product, error) {
	query := "select * from product "
	/*
		where {filter, pencarian}
		non where {urutkan}
	*/

	if filter != "" || pencarian != "" {

		query += " where "

		if filter != "" && pencarian != "" {
			// semua
			query += "judul_product ilike '%" + pencarian + "%'"
			query += " and "
			query += "(" + filter + " >= " + dari + " and " + filter + " <= " + hingga + ")"

		} else {
			// salah satu
			if pencarian != "" {
				query += "judul_product ilike '%" + pencarian + "%'"
			}
			if filter != "" {
				query += "(" + filter + " >= " + dari + " and " + filter + " <= " + hingga + ")"
			}
		}
	}

	if urutkan != "" {
		if berdasar != "" {
			// sesuai filter
			query += "order by " + urutkan + " " + berdasar
		} else {
			// default ASC
			query += "order by " + urutkan + " ASC"
		}
	}

	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []Product{} //function return data

	for rows.Next() {

		product := Product{} //kontainer for all komentar
		//var isikomentar string
		err = rows.Scan(&product.IDProduct, &product.FotoProduct, &product.JudulProduct, &product.DeskripsiProduct, &product.Terjual, &product.Disukai, &product.Harga)
		//komentar.IsiKomentar = template.HTMLEscapeString(isikomentar)
		if err != nil {
			log.Println(err.Error())
		}
		data = append(data, product)
		err = rows.Err()
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	return data, nil
}
