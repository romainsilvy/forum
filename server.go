package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	databaseTools "test/dataBase"

	// "forum/dataBase/database"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/gorilla/sessions"
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
func getIncrisption(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")
		http.SetCookie(w, &http.Cookie{
			Name:  "my-cookie",
			Value: "some value",
			Path:  "/",
		})
		fmt.Println("COOKIE CREATED")

		if r.FormValue("inscriptionPseudo") != "" {
			inscriptionPseudo := r.FormValue("inscriptionPseudo")
			inscriptionEmail := r.FormValue("inscriptionEmail")
			inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
			inscriptionPassword := r.FormValue("inscriptionPassword")
			inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

			if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
				databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "test")
			}
		}
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

		// fmt.Println(r.FormValue("inscriptionPseudo"))
		// fmt.Println(r.FormValue("inscriptionEmail"))
		// fmt.Println(r.FormValue("inscriptionEmailConfirm"))
		// fmt.Println(r.FormValue("inscriptionPassword"))
		// fmt.Println(r.FormValue("inscriptionPasswordConfirm"))
	})
}

// func getUsers(oneUser databaseTools.User, tabUser []databaseTools.User) {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		variable, _ := template.ParseFiles("index.html")
// 		database, _ := sql.Open("sqlite3", "dataBase/forum.db")
// 		rows, _ := database.Query("select * from User")
// 		result := tabUser
// 		for rows.Next() {
// 			item := oneUser
// 			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
// 			if err2 != nil {
// 				panic(err2)
// 			}
// 			result = append(result, item)
// 		}
// 		variable.Execute(w, result)
// 	})
// }

func handleProfil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("profil.html")
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

func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: "some value",
		Path:  "/",
	})
	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
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

func main() {
	databaseOpened, _ := sql.Open("sqlite3", "dataBase/forum.db")
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleConnexion()
	handleInscription()
	getIncrisption(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	// mux.HandleFunc("/createcookie", CreateCookie)
	// http.ListenAndServe(":8080", mux)
	// getUsers(databaseTools.User{}, []databaseTools.User{})
	runServer()
}

// func secret(w http.ResponseWriter, r *http.Request) {
//     session, _ := store.Get(r, "cookie-name")

//     // Check if user is authenticated
//     if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
//         http.Error(w, "Forbidden", http.StatusForbidden)
//         return
//     }

//     // Print secret message
//     fmt.Fprintln(w, "The cake is a lie!")
// }

// func login(w http.ResponseWriter, r *http.Request) {
//     session, _ := store.Get(r, "cookie-name")

// 	//CHECK PASSWORD AND EMAIL
//     session.Values["authenticated"] = true
//     session.Save(r, w)
// }

// func logout(w http.ResponseWriter, r *http.Request) {
//     session, _ := store.Get(r, "cookie-name")

//     // Revoke users authentication
//     session.Values["authenticated"] = false
//     session.Save(r, w)
// }

// func mainConnexion() {
//     http.HandleFunc("/secret", secret)
//     http.HandleFunc("/login", login)
//     http.HandleFunc("/logout", logout)

//     http.ListenAndServe(":8080", nil)
// }

//LIRE BDD ET ENVOYER

// for rows.Next() {
// 	item := oneUser
// 	err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
// 	if err2 != nil {
// 		panic(err2)
// 	}
// 	result = append(result, item)
// }
// variable.Execute(w, result)
