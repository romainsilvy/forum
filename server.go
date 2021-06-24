package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	AccountManagement "tools/AccountManagement"
	displayTools "tools/Display"
	interractionsTools "tools/Interractions"
	databaseTools "tools/dataBase"

	_ "github.com/mattn/go-sqlite3"
)

//handleAccueil is the handlefunc for the main page
func handleAccueil(database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var dataToSend []databaseTools.ThreadData
		variable, _ := template.ParseFiles("index.html")

		AccountManagement.Connexion(w, r, database)
		AccountManagement.Inscription(r, database)
		dataToSend = displayTools.RetrieveAccueil(dataToSend, w, database)

		variable.Execute(w, dataToSend)
	})
}

// Run all the necessary functions for the life of the serv
func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(db)
	AccountManagement.HandleProfil(databaseTools.User{}, db)
	interractionsTools.FetchLike(db)
	interractionsTools.FetchThread(db)
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Run the serv
func main() {
	databaseTools.InitDatabase("dataBase/forum.db")
	db, _ := sql.Open("sqlite3", "dataBase/forum.db")
	handleAll(db)
	runServer()
}
