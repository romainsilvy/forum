package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id_user   int
	User_name string
	Password  string
	Email     string
}

// get Data from the sqlite database and print them int the html page
func getUsers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")
		database, _ := sql.Open("sqlite3", "./forum.db")
		rows, _ := database.Query("select * from User")
		var result []User
		for rows.Next() {
			item := User{}
			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, item)
		}
		variable.Execute(w, result)
	})
}

func main() {
	getUsers()
}
