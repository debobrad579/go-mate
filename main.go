package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app/index.html")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" {
			path = "/index"
		}

		filePath := filepath.Join("views", strings.TrimPrefix(path, "/")+".html")

		http.ServeFile(w, r, filePath)
	})

	http.ListenAndServe(":3000", mux)
}
