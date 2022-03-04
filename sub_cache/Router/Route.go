

package Router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"sub_cache/HttpProcessing"
	"sub_cache/Models"
)

type IpgDataBase interface {
	GetOrders(context.Context, string) (Models.OrderInfo, error)
}

type Postgres struct {
	DB IpgDataBase
}

func Route(pg Postgres, cache *lru.Cache) *chi.Mux{
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/cache/{id}", func(w http.ResponseWriter, r *http.Request){
		id := chi.URLParam(r, "id")

		var order interface{}
		var ok bool
		var errGetOrder error

		fmt.Println("get from cache")
		order, ok = cache.Get(id)
		if !ok {
			fmt.Println("get from db")
			order, errGetOrder = pg.DB.GetOrders(r.Context(), id)
			if errGetOrder != nil {
				HttpProcessing.ErrorHandler(w, errGetOrder, "get order",
					"server error", http.StatusInternalServerError)
				return
			}

			cache.Add(order.(Models.OrderInfo).ID, order)
		}

		resp, errMarshal := json.Marshal(order.(Models.OrderInfo))
		if errMarshal != nil {
			HttpProcessing.ErrorHandler(w, errMarshal, "marshal order",
				"server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application-json")
		w.Write(resp)
	})

	return router
}
