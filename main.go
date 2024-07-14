package main

import (
	"html/template"
	"log"
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

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL.Path)
			f(w, r)
		}
	}
}

func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			f(w, r)
		}
	}
}

func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}

func HomeController(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
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
	tmpl.Execute(w, data)
}

func FormController(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
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
}

func main() {

	fs := http.FileServer(http.Dir("assets/"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if tmpl, err := template.ParseFiles("views/template.html"); err != nil {
		panic("There is not html")
	} else {
		http.HandleFunc("/", Chain(func(w http.ResponseWriter, r *http.Request) {
			HomeController(w, r, tmpl)
		}, Method("GET"), Logging()))
	}

	if tmpl, err := template.ParseFiles("views/forms.html"); err != nil {
		panic("There is not html")
	} else {
		http.HandleFunc("/form", Chain(func(w http.ResponseWriter, r *http.Request) {
			FormController(w, r, tmpl)
		}, Method("GET"), Logging()))
	}

	http.ListenAndServe(":3000", nil)
}
