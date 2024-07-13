package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	fs := http.FileServer(http.Dir("assets/"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if tmpl, err := template.ParseFiles("views/template.html"); err != nil {
		panic("There is not html")
	} else {
		data := TodoPageData{
			PageTitle: "List Of Todo",
			Todos: []Todo{
				{
					Title: "Go to the store",
					Done:  false,
				},
				{
					Title: "Make some screencast",
					Done:  true,
				},
				{
					Title: "Get some intermezzo",
					Done:  false,
				},
			},
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl.Execute(w, data)
		})
	}

	http.ListenAndServe(":3000", nil)
}
