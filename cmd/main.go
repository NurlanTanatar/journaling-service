package main

import (
	"log"
	"net/http"
)

// TODO: save svg as asset
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
	err = parseTemplates()
	if err != nil {
		log.Panic(err)
	}
	r := http.NewServeMux()
	r.Handle("GET /static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("GET /", handleGetIncidents)
	r.HandleFunc("POST /tasks", handleCreateIncident)
	r.HandleFunc("PUT /tasks/{id}/toggle", handleToggleIncident)
	r.HandleFunc("DELETE /tasks/{id}", handleDeleteIncident)
	r.HandleFunc("GET /tasks/{id}/edit", handleEditIncident)
	r.HandleFunc("PUT /tasks/{id}", handleUpdateIncident)
	r.HandleFunc("PUT /tasks", handleOrderIncidents)
	http.ListenAndServe(":3000", r)
}
