package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/poc/service"
)

//Routes module endpoints
func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("/financials/*", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("testing")
		render.Status(r, http.StatusOK)
		response := service.Response{}
		response.Code = http.StatusOK
		response.Data = "Financials Service - Response"
		render.JSON(rw, r, response)
	})
	return r
}
