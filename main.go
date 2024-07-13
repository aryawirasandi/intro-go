package main

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	conn := "your_conn_config"

	db, err := sql.Open("mysql", conn)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// welcoming screen
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	usersRouter := r.PathPrefix("/users").Subrouter()

	usersRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		type User struct {
			id        int
			username  string
			password  string
			createdAt string
		}

		rows, err := db.Query("SELECT * FROM users")

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var users []User

		for rows.Next() {
			var u User

			if err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		if len(users) == 0 {
			fmt.Fprint(w, "user is empty")
		}

		for _, user := range users {
			fmt.Fprintf(w, "Username = %v", user.username)
		}
	})
	usersRouter.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		// vars := mux.Vars(r);

		username := r.FormValue("username")
		password := r.FormValue("password")
		createdAt := time.Now()

		result, err := db.Exec("INSERT INTO users (username, password, created_at) VALUES (?,?,?)", username, password, createdAt)

		if err != nil {
			fmt.Fprint(w, err)
		}

		if _, err := result.LastInsertId(); err != nil {
			fmt.Fprint(w, err)
		} else {
			fmt.Fprint(w, "Successfully Created user")
		}
	}).Methods("POST")

	usersRouter.HandleFunc("/show/{id}", func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT * FROM users WHERE id = ?
		`

		v := mux.Vars(r)

		idUser := v["id"]

		var (
			id        int
			username  string
			password  string
			createdAt string
		)

		err := db.QueryRow(query, idUser).Scan(&id, &username, &password, &createdAt)

		if err != nil {
			fmt.Fprint(w, err)
		}

		fmt.Fprintf(w, "\n"+fmt.Sprint(id))
		fmt.Fprintf(w, "\n"+username)
		fmt.Fprintf(w, "\n"+password)
	})

	usersRouter.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		query := `DELETE FROM users WHERE id = ?`

		v := mux.Vars(r)

		idUser := v["id"]

		_, err := db.Exec(query, idUser)

		if err != nil {
			fmt.Fprint(w, err)
		}

		fmt.Fprint(w, "Successfully Delete User")

	})
	http.ListenAndServe(":3000", r)
}
