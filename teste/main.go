package main

import (
	"net/http"

	"github.com/andreluzz/poc/modules/shared"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	router := chi.NewRouter()
	router.Use(
		middleware.Heartbeat("/ping"),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Mount("/api/v1", routes())

	http.ListenAndServe(":3000", router)
}

func routes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("/core/*", func(rw http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		response := shared.Response{}
		response.Code = http.StatusOK
		response.Data = "Core Service - Response"
		render.JSON(rw, r, response)
	})
	return r
}
