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

	if tmpl, err := template.ParseFiles("views/forms.html"); err != nil {
		panic("There is not html")
	} else {
		// data := TodoPageData{
		// 	PageTitle: "List Of Todo",
		// 	Todos: []Todo{
		// 		{
		// 			Title: "Go to the store",
		// 			Done:  false,
		// 		},
		// 		{
		// 			Title: "Make some screencast",
		// 			Done:  true,
		// 		},
		// 		{
		// 			Title: "Get some intermezzo",
		// 			Done:  false,
		// 		},
		// 	},
		// }
		http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				tmpl.Execute(w, nil)
				return
			}

			type ContatDetails struct {
				Email   string
				Subject string
				Message string
			}

			details := ContatDetails{
				Email:   r.FormValue("email"),
				Subject: r.FormValue("subject"),
				Message: r.FormValue("message"),
			}

			result := details

			tmpl.Execute(w, struct {
				Success bool
				Result  ContatDetails
			}{true, result})

		})
	}

	http.ListenAndServe(":3000", nil)
}
