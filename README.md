# atlas
Atlas is a simple image/video file server that works with vanilla `net/http` server.
It supports "range headers" and image compression.

## Example
Atlas will serve all files in all subfolders for the provided root.

See `examples/basic_usage` for an example
```go
package main

import (
	"fmt"
	"github.com/k0rean-rand0m/atlas"
	"net/http"
)

func main() {
	http.HandleFunc("/", atlas.ServeMedia("examples/basic_usage/media"))
	fmt.Println("http://localhost:8008/rick/roll.webp")
	http.ListenAndServe(":8008", nil)
}
```