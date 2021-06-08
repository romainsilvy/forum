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

		variable, _ := template.ParseFiles("inscription.html")
		result := 3
		variable.Execute(w, result)
	})
}

func handleVoirprofil() {
	http.HandleFunc("/voirprofil", func(w http.ResponseWriter, r *http.Request) {

		variable, _ := template.ParseFiles("voir-profil.html")
		result := 3
		variable.Execute(w, result)
	})
}

func handleProfil(oneUser database.User, tabUser []database.User) {
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("profil.html")
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

func main() {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// database.InsertIntoUsers("Amaury", "mail", "mdp", "image")
	getUsers(database.User{}, []database.User{})
	handleInscription()
	handleConnexion()
	handleVoirprofil()
	handleProfil(database.User{}, []database.User{})
	// cssFile()
	// pictureFile()
	runServer()
}
