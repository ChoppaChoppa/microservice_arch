package Route

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"sub_db/HttpProcessing"
	"sub_db/Models"
)

type IPgDataBase interface {
	Get(ctx context.Context, id string) (Models.OrderInfo, error)
}

type DataBase struct {
	DB IPgDataBase
}

func Router(pg DataBase) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		order, errGet := pg.DB.Get(r.Context(), id)
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

	return router
}
