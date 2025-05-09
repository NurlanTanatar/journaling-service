package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	layout = "2006-01-02T15:04"
)

func handleGetIncidents(w http.ResponseWriter, _ *http.Request) {
	items, err := fetchIncidents()
	if err != nil {
		log.Printf("error fetching tasks: %v", err)
		return
	}
	count, err := fetchCount()
	if err != nil {
		log.Printf("error fetching count: %v", err)

	}
	completedCount, err := fetchCompletedCount()
	if err != nil {
		log.Printf("error fetching completedCount: %v", err)
	}
	data := Incidents{
		Items:          items,
		Count:          count,
		CompletedCount: completedCount,
	}
	tmpl.ExecuteTemplate(w, "Base", data)
}

func handleCreateIncident(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	criticality, err := strconv.Atoi(r.FormValue("criticality"))
	if err != nil {
		log.Printf("error parsing and getting criticalty %v", err)
	}
	dateStart, err := time.Parse(layout, r.FormValue("date-start"))
	if err != nil {
		log.Printf("error parsing and getting dateStart %v", err)
	}
	dateEnd, err := time.Parse(layout, r.FormValue("date-end"))
	if err != nil {
		log.Printf("error parsing and getting dateEnd %v", err)
	}

	if title == "" {
		tmpl.ExecuteTemplate(w, "Form", nil)
		return
	}

	item, err := insertIncident(title, criticality, dateStart, dateEnd)
	if err != nil {
		log.Printf("error inserting task: %v", err)
		return
	}
	count, err := fetchCount()
	if err != nil {
		log.Printf("error fetching count: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	tmpl.ExecuteTemplate(w, "Form", nil)
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": item, "SwapOOB": true})
	tmpl.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})
}

func handleToggleIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("error parsing id into int %v", err)
		return
	}
	_, err = toggleIncident(id)
	if err != nil {
		log.Printf("error toggling task %v", err)
		return
	}
	completedCount, err := fetchCompletedCount()
	if err != nil {
		log.Printf("error fetching completed count %v", err)
		return
	}
	tmpl.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": completedCount, "SwapOOB": true})
}

func handleDeleteIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("error parsing id into int %v", err)
		return
	}
	err = deleteIncident(r.Context(), id)
	if err != nil {
		log.Printf("error deleting task %v", err)
		return
	}
	count, err := fetchCount()
	if err != nil {
		log.Printf("error fetching count %v", err)
	}
	completedCount, err := fetchCompletedCount()
	if err != nil {
		log.Printf("error fetching completed count %v", err)
	}
	tmpl.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})
	tmpl.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": completedCount, "SwapOOB": true})
}

func handleEditIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("error parsing id into int %v", err)
		return
	}
	task, err := fetchIncident(id)
	if err != nil {
		log.Printf("error fetching task with id %v", err)
		return
	}
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": task, "Editing": true})
}

func handleUpdateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("error parsing id into int %v", err)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		return
	}
	task, err := updateIncident(id, title)
	if err != nil {
		log.Printf("error updating task with id %v", err)
		return
	}
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": task})
}

func handleOrderIncidents(_ http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form %v", err)
		return
	}
	values := []int{}
	for k, v := range r.Form {
		if k == "item" {
			for _, v := range v {
				value, err := strconv.Atoi(v)
				if err != nil {
					log.Printf("error parsing id into int %v", err)
					return
				}
				values = append(values, value)
			}
		}
	}
	err = orderIncidents(r.Context(), values)
	if err != nil {
		log.Printf("error ordering tasks %v", err)
		return
	}
}
