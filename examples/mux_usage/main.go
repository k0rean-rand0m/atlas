package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/k0rean-rand0m/atlas"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc(
		"/m/{rest:.*}",
		atlas.Handler("/m", "examples/mux_usage/media"),
	).Methods("GET")

	fmt.Println("http://localhost:8008/m/rick/roll.webp")
	http.Handle("/", r)
	http.ListenAndServe(":8008", nil)
}
