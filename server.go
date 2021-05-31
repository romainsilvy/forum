package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id_user   int
	User_name string
	Password  string
	Email     string
}

// serveFile makes files available for the website
func cssFile() {
	cssServer := http.FileServer(http.Dir("./data/css"))
	http.Handle("/global.css", cssServer)
	http.Handle("/accueil.css", cssServer)
	http.Handle("/accueil-droit.css", cssServer)
}

func pictureFile() {
	pictureServer := http.FileServer(http.Dir("./data/images"))
	http.Handle("/logo.png", pictureServer)
	http.Handle("/profile-picture.png", pictureServer)
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
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
	cssFile()
	pictureFile()
	runServer()
}
