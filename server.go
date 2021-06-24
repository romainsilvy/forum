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

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//handleAccueil is the handlefunc for the main page
func handleAccueil(database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var dataToSend []databaseTools.ThreadToSend
		variable, _ := template.ParseFiles("index.html")

		//add thread variables
		title := r.FormValue("threadTitle")
		content := r.FormValue("créa_thread")
		submitButton := r.FormValue("submitthread")

		//categories variables
		// inputCatCham := r.FormValue("CHAMEAU")
		// inputCatDrom := r.FormValue("DROMADAIRE")
		// inputCatLama := r.FormValue("LAMA")

		//supp variable
		deleteButton := r.FormValue("suppr")

		//session cookie
		sessionCookieAuth, _ := store.Get(r, "auth")

		// inputSearchBar := r.FormValue("searchWord")
		inputCatThread := r.FormValue("drone")

		if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] == true) {
			interractionsTools.AddThread(sessionCookieAuth, title, content, inputCatThread, database)
		}

		if (deleteButton != "") && (sessionCookieAuth.Values["authenticated"] == true) {
			interractionsTools.SuppThread(sessionCookieAuth, database)
		}

		AccountManagement.Connexion(w, r, database)
		AccountManagement.Inscription(r, database)

		// if inputCatCham != "" {
		// 	displayTools.DisplayCategory(inputCatCham, dataToSend, w, database)
		// 	variable.Execute(w, dataToSend)
		// } else if inputCatDrom != "" {
		// 	displayTools.DisplayCategory(inputCatDrom, dataToSend, w, database)
		// 	variable.Execute(w, dataToSend)
		// } else if inputCatLama != "" {
		// 	displayTools.DisplayCategory(inputCatLama, dataToSend, w, database)
		// 	variable.Execute(w, dataToSend)
		// } else if inputSearchBar != "" {
		// 	displayTools.DisplaySearchResult(inputSearchBar, dataToSend, w, database)
		// 	variable.Execute(w, dataToSend)
		// } else {
		displayTools.DisplayAccueil(dataToSend, w, database)
		variable.Execute(w, dataToSend)
		// }
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
