package router

import (
	"main/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/dashboard", controller.GetMessage).Methods("GET")
	r.HandleFunc("/dashboard", controller.ConvertMessagetoTXT).Methods("POST")

	return r
}
