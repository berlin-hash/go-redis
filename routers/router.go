package routers

import (
	"go-redis/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.GetProductUsingID).Methods("GET")
	return router
}
