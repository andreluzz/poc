package routes

import (
	"net/http"

	"github.com/andreluzz/poc/modules/shared"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//Routes module endpoints
func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("/financials/*", func(rw http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		response := shared.Response{}
		response.Code = http.StatusOK
		response.Data = "Financials Service - Response"
		render.JSON(rw, r, response)
	})
	return r
}
