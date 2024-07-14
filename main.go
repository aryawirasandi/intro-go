package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func decode(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	fmt.Fprintf(w, "%s %s is %d years old!", user.Firstname, user.Lastname, user.Age)
}

func encode(w http.ResponseWriter, r *http.Request) {
	john := User{
		Firstname: "john",
		Lastname:  "doe",
		Age:       42,
	}

	json.NewEncoder(w).Encode(john)
}

func main() {
	http.HandleFunc("/decode", decode)
	http.HandleFunc("/encode", encode)

	http.ListenAndServe(":3000", nil)
}
