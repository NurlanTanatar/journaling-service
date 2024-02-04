package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	err := openDB()
	if err != nil {
		log.Panic(err)
	}
	defer closeDB()
	err = setupDB()
	if err != nil {
		log.Panic(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("").ParseFiles("templates/index.html")
		tmpl.ExecuteTemplate(w, "Base", nil)
	})
	http.ListenAndServe(":3000", r)
}
