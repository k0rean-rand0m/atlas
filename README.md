# atlas
Atlas is a simple image/video file server that works with vanilla `net/http` server.
It supports "range headers" and image compression.

## Examples
Atlas will serve all files in all subfolders for the provided root.

### With net/http
See `examples/basic_usage`
```go
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
```

### With gorilla/mux
See `examples/mux_usage`
```go
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

```