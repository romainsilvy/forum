package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	database "test/dataBase"

	// "forum/dataBase/database"
	_ "github.com/mattn/go-sqlite3"
)

// type User struct {
// 	Id_user   int
// 	User_name string
// 	Password  string
// 	Email     string
// }

// serveFile makes files available for the website
func cssFile() {
	cssServer := http.FileServer(http.Dir("./data/css"))
	http.Handle("/global.css", cssServer)
	http.Handle("/accueil.css", cssServer)
	http.Handle("/accueil-droit.css", cssServer)
	http.Handle("/monprofil.css", cssServer)
	http.Handle("/accueil-gauche.css", cssServer)
	http.Handle("/connexion.css", cssServer)
	http.Handle("/inscription.css", cssServer)
}

func pictureFile() {
	pictureServer := http.FileServer(http.Dir("./data/images"))
	http.Handle("/logo.png", pictureServer)
	http.Handle("/profile-picture.png", pictureServer)
	http.Handle("/follow.png", pictureServer)
	http.Handle("/home.png", pictureServer)
	http.Handle("/thread.png", pictureServer)
	http.Handle("/like.png", pictureServer)
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// get Data from the sqlite database and print them int the html page
func getUsers(oneUser database.User, tabUser []database.User) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")
		database, _ := sql.Open("sqlite3", "dataBase/forum.db")
		rows, _ := database.Query("select * from User")
		result := tabUser
		for rows.Next() {
			item := oneUser
			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, item)
		}
		variable.Execute(w, result)
	})
}

func handleConnexion() {
	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {

		variable, _ := template.ParseFiles("connexion.html")
		result := 3
		variable.Execute(w, result)
	})
}

func handleInscription() {
	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {

		variable, _ := template.ParseFiles("connexion.html")
		result := 3
		variable.Execute(w, result)
	})
}

func handleProfil() {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		pseudoPath := r.URL.Path[7:]
		variable, _ := template.ParseFiles("profil.html")
		variable.Execute(w, pseudoPath)
	})
}

func main() {
	// database.InsertIntoUsers("Amaury", "mail", "mdp", "image")
	getUsers(database.User{}, []database.User{})
	handleInscription()
	handleConnexion()
	handleProfil()
	cssFile()
	pictureFile()
	runServer()
}
