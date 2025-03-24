package main

import (
	"fmt"
	"github.com/k0rean-rand0m/atlas"
	"net/http"
)

func main() {
	http.HandleFunc("/", atlas.Handler("/static", "examples/basic_usage/media"))
	fmt.Println("http://localhost:8008/static/rick/roll.webp")
	http.ListenAndServe(":8008", nil)
}
