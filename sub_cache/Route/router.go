package Route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Router(){
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request){
		id := chi.URLParam(r, "id")

		//TODO обработка запрсса и отпрвка ответа
	})
}
