

package Router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
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

func Route(cache *lru.Cache) *chi.Mux{
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/cache/{id}", func(w http.ResponseWriter, r *http.Request){
		id := chi.URLParam(r, "id")

		order, ok := cache.Get(id)
		if !ok {
			fmt.Println("get from db")

			orderResp, errResp := http.Get("http://127.0.0.1:3000/sub_db/get/" + id)
			if errResp != nil {
				HttpProcessing.ErrorHandler(w, errResp, "http request",
					"server error", http.StatusInternalServerError)
				return
			}
			body, errGetBody := ioutil.ReadAll(orderResp.Body)
			if errGetBody != nil {
				HttpProcessing.ErrorHandler(w, errGetBody, "get body",
					"server error", http.StatusInternalServerError)
				return
			}

			var modelOrder Models.OrderInfo
			if errGetOrder := json.Unmarshal(body, &modelOrder); errGetOrder != nil {
				HttpProcessing.ErrorHandler(w, errGetOrder, "unmarshal body",
					"server error", http.StatusInternalServerError)
				return
			}

			cache.Add(modelOrder.ID, modelOrder)
			sendResp(w, modelOrder)
			return
		}

		fmt.Println("get from cache")
		sendResp(w, order.(Models.OrderInfo))
	})

	return router
}

func sendResp(w http.ResponseWriter, order Models.OrderInfo){
	resp, errMarshal := json.Marshal(order)
	if errMarshal != nil {
		HttpProcessing.ErrorHandler(w, errMarshal, "marshal order",
			"server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.Write(resp)
}