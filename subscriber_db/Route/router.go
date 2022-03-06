package Route

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"net/http"

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

		route.Get("/lasts", func(w http.ResponseWriter, r *http.Request) {
			type request struct {
				count int
			}
			var reqStruct request
			body, errGetBody := ioutil.ReadAll(r.Body)
			if errGetBody != nil {
				HttpPorcessing.HttpError(w, errGetBody, "read body",
					"bad request", http.StatusBadRequest)
				return
			}
			if errUnmarshalBody := json.Unmarshal(body, &reqStruct); errUnmarshalBody != nil {
				HttpPorcessing.HttpError(w, errUnmarshalBody, "unmarshal body",
					"bad request", http.StatusBadRequest)
				return
			}

			orders, errGet := pg.DB.GetLasts(r.Context(), reqStruct.count)
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
