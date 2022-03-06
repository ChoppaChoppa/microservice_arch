package Route

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"

	"sub_db/HttpProcessing"
	"sub_db/Models"
)

type IPgDataBase interface {
	GetOrder(context.Context, string) (Models.OrderInfo, error)
	GetLasts(context.Context, int)    ([]Models.OrderInfo, error)
}

type DataBase struct {
	DB IPgDataBase
}

func Router(pg DataBase) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/sub_db/get", func(route chi.Router) {
		route.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			order, errGet := pg.DB.GetOrder(r.Context(), id)
			if errGet != nil {
				HttpPorcessing.HttpError(w, errGet, "err get",
					"server error", http.StatusInternalServerError)
				return
			}

			resp, errMarshal := json.Marshal(order)
			if errMarshal != nil {
				HttpPorcessing.HttpError(w, errMarshal, "err marshal",
					"server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			w.Write(resp)
		})

		route.Get("/lasts/{count}", func(w http.ResponseWriter, r *http.Request) {
			param := chi.URLParam(r, "count")

			count, errIsDigit := strconv.Atoi(param)
			if errIsDigit != nil {
				HttpPorcessing.HttpError(w, errIsDigit, "param is not digit",
					"param is not digit", http.StatusBadRequest)
				return
			}

			orders, errGet := pg.DB.GetLasts(r.Context(), count)
			if errGet != nil {
				HttpPorcessing.HttpError(w, errGet, "get last orders",
					"server error", http.StatusInternalServerError)
				return
			}

			resp, errMarshal := json.Marshal(orders)
			if errMarshal != nil {
				HttpPorcessing.HttpError(w, errGet, "marshal orders",
					"server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			w.Write(resp)
		})
	})

	return router
}
