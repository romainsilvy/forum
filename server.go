package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	databaseTools "test/dataBase"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//inscription manage the inscription form
func inscription(r *http.Request) {
	inscriptionPseudo := r.FormValue("inscriptionPseudo")
	if inscriptionPseudo != "" {
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "img")
		}
	}
}

//connexion manage the connexion form
func connexion(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.FormValue("connect") != "" {
		connexionUser := r.FormValue("connexionUser")
		connexionPassword := r.FormValue("connexionPassword")
		checkIfExist := databaseTools.CheckIfExist(database, "password", "User", "user_name", connexionUser)

		if checkIfExist {
			userPassword := databaseTools.SingleRowQuerry(database, "password", "User", "user_name", connexionUser)
			if userPassword != "notExist" {
				if userPassword == connexionPassword {
					session, _ := store.Get(r, "auth")
					session.Values["authenticated"] = true
					session.Values["user"] = connexionUser
					session.Save(r, w)
					fmt.Println("utilisateur connecté")
				} else {
					fmt.Println("mot de passe incorrect")
				}
			}
		} else {
			fmt.Println("pseudo incorrect")
		}
	}
}

//handleAccueil is the handlefunc for the main page
func handleAccueil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")
		inscription(r)
		connexion(w, r, database)
		variable.Execute(w, tabUser)
	})
}

func changePassword(r *http.Request, userPassword string, userName string, db *sql.DB) {
	changePasswordButton := r.FormValue("changePasswordButton")
	if changePasswordButton != "" {
		oldPassword := r.FormValue("oldPassword")
		if oldPassword == userPassword {
			newPassword := r.FormValue("newPassword")
			newPasswordConfirm := r.FormValue("newPasswordConfirm")
			if newPassword == newPasswordConfirm {
				// databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "img")
				databaseTools.UpdateValue(db, "User", "password", newPassword, "user_name", userName)
				fmt.Println("mot de passe changé ")
			} else {
				fmt.Println("confirmation de mdp pas bonne")
			}
		} else {
			fmt.Println("ancien mdp pas bon")
		}
	}
}

func changeEmail(r *http.Request, userPassword string, userName string, db *sql.DB) {
	changeEmailButton := r.FormValue("changeEmailButton")
	if changeEmailButton != "" {
		newEmail := r.FormValue("newEmail")
		newEmailConfirm := r.FormValue("newEmailConfirm")

		if newEmail == newEmailConfirm {
			password := r.FormValue("password")
			if password == userPassword {
				databaseTools.UpdateValue(db, "User", "email", newEmail, "user_name", userName)
				fmt.Println("change email")
			} else {
				fmt.Println("mauvais mot de passe ")

			}
		} else {
			fmt.Println("email differents")
		}
	}
}

//handleProfil is the handlefunc for the profil page
func handleProfil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("profil.html")

		session, _ := store.Get(r, "auth")
		username := session.Values["user"].(string)
		oneUser.User_name = username
		oneUser.Email = databaseTools.SingleRowQuerry(database, "email", "User", "user_name", oneUser.User_name)
		oneUser.Password = databaseTools.SingleRowQuerry(database, "password", "User", "user_name", oneUser.User_name)

		changePassword(r, oneUser.Password, oneUser.User_name, database)
		changeEmail(r, oneUser.Password, oneUser.User_name, database)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		variable.Execute(w, oneUser)
	})
}

func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(databaseTools.User{}, []databaseTools.User{}, db)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, db)
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, _ := sql.Open("sqlite3", "dataBase/forum.db")
	handleAll(db)
	runServer()
}
