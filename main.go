package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	mux.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app/index.html")
	})

	mux.Handle("/", http.FileServer(http.Dir("views")))

	http.ListenAndServe(":3000", mux)
}
