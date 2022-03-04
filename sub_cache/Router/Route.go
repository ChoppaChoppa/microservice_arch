package Router

import(
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Route(){
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("cache/{id}", func(w http.ResponseWriter, r *http.Request){
		id := chi.URLParam(r, "id")


	})
}
