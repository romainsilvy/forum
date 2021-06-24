package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	AccountManagement "tools/AccountManagement"
	interractionsTools "tools/Interractions"
	databaseTools "tools/dataBase"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

// Display the category and append values from sql Thread table
func displayCategory(inputCatChoisie string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveCategoryRows(db, inputCatChoisie)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	variable.Execute(w, dataToSend)
}

// Display the research from the searchBar and append values from sql Thread table
func displaySearchResult(inputSearchBar string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveSearchRows(db, inputSearchBar)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	variable.Execute(w, dataToSend)
}

// Display all the threads on the home page and append the values in sql Thread Table
func displayAccueil(dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveAccueilRows(db)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	variable.Execute(w, dataToSend)
}

//handleAccueil is the handlefunc for the main page
func handleAccueil(database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var dataToSend []databaseTools.ThreadData
		variable, _ := template.ParseFiles("index.html")

		//add thread variables
		title := r.FormValue("threadTitle")
		content := r.FormValue("créa_thread")
		submitButton := r.FormValue("submitthread")

		//categories variables
		inputCatCham := r.FormValue("CHAMEAU")
		inputCatDrom := r.FormValue("DROMADAIRE")
		inputCatLama := r.FormValue("LAMA")

		//supp variable
		deleteButton := r.FormValue("suppr")

		//session cookie
		sessionCookieAuth, _ := store.Get(r, "auth")

		inputSearchBar := r.FormValue("searchWord")
		inputCatThread := r.FormValue("drone")

		if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] == true) {
			interractionsTools.AddThread(sessionCookieAuth, title, content, inputCatThread, database)
		}

		if (deleteButton != "") && (sessionCookieAuth.Values["authenticated"] == true) {
			interractionsTools.SuppThread(sessionCookieAuth, database)
		}

		AccountManagement.Connexion(w, r, database)
		AccountManagement.Inscription(r, database)

		if inputCatCham != "" {
			displayCategory(inputCatCham, dataToSend, variable, w, database)
		} else if inputCatDrom != "" {
			displayCategory(inputCatDrom, dataToSend, variable, w, database)
		} else if inputCatLama != "" {
			displayCategory(inputCatLama, dataToSend, variable, w, database)
		} else if inputSearchBar != "" {
			displaySearchResult(inputSearchBar, dataToSend, variable, w, database)
		} else {
			displayAccueil(dataToSend, variable, w, database)
		}
	})
}

// Run all the necessary functions for the life of the serv
func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(db)
	AccountManagement.HandleProfil(databaseTools.User{}, db)
	interractionsTools.FetchLike(db)
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
