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

    // "github.com/gorilla/sessions"
)

// type User struct {
// 	Id_user   int
// 	User_name string
// 	Password  string
// 	Email     string
// }



func InsertIntoDb(user_name string, password string) {
	db := database.InitDatabase("dataBase/forum.db")
	_, err := db.Exec(`INSERT INTO User (user_name, password) VALUES (?, ?)`, user_name, password)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
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


		http.SetCookie(w, &http.Cookie{
			Name:  "my-cookie", 
			Value: "some value",  
			Path: "/",
   		})
   		fmt.Println("COOKIE CREATED")

		variable, _ := template.ParseFiles("index.html")
		database, _ := sql.Open("sqlite3", "dataBase/forum.db")
		rows, _ := database.Query("select * from User")
		result := tabUser

		nameUser := r.FormValue("inputNameUser")
		passwordUser := r.FormValue("inputPasswordUser")
		InsertIntoDb(nameUser, passwordUser)
		fmt.Println(nameUser)
		fmt.Println(passwordUser)


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

func handleProfil(oneUser database.User, tabUser []database.User) {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
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








func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
			 Name:  "my-cookie", 
			 Value: "some value",  
			 Path: "/",
	})
	fmt.Fprintln(w, "COOKIE WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(w, "in chrome go to: dev tools / application / cookies")
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








func main() {
	// mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// database.InsertIntoUsers("Amaury", "mail", "mdp", "image")
	getUsers(database.User{}, []database.User{})
	handleInscription()
	handleConnexion()
	handleProfil(database.User{}, []database.User{})
	// mux.HandleFunc("/createcookie", CreateCookie)
	// http.ListenAndServe(":8080", mux)
	runServer()
}
