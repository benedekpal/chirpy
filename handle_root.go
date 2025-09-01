package main

import (
	"net/http"
)

func registerRoot(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("./static")))
}
