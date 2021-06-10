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

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

//inscription manage the inscription popup 
func inscription(r *http.Request)  {
	if r.FormValue("inscriptionPseudo") != "" {
		inscriptionPseudo := r.FormValue("inscriptionPseudo")
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "img")
		}
	}
}

//connexion manage the connexion popup 
func connexion(w http.ResponseWriter, r *http.Request, database *sql.DB)  {
	if r.FormValue("connect") != "" {
		connexionUser := r.FormValue("connexionUser")
		connexionPassword := r.FormValue("connexionPassword")
		checkIfExist := checkIfExist(database, "password", "User", "user_name", connexionUser)

		if checkIfExist {
			userPassword := singleRowQuerry(database, "password", "User", "user_name", connexionUser)
			if userPassword != "notExist" {
				if userPassword == connexionPassword {
					session, _ := store.Get(r, "auth")
					session.Values["authenticated"] = true
					session.Values["user"] = connexionUser
					session.Save(r, w)
					fmt.Println("tu es co chacal")
				} else {
					fmt.Println("mauvais mdp chacal")
				}
			}
		} else {
			fmt.Println("mauvais pseudo chacal")
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

//handleProfil is the handlefunc for the profil page 
func handleProfil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("profil.html")
		
		session, _ := store.Get(r, "auth")
		username := session.Values["user"].(string)
		oneUser.User_name = username
		oneUser.Email = singleRowQuerry(database, "email", "User", "user_name", oneUser.User_name)
		
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		variable.Execute(w, oneUser)
	})
}

//singleRowQuerry retrieve a value in the db with a where comparator
func singleRowQuerry(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) string {
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}
	// defer stmt.Close()
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		return "notExist"
	}
	return toReturn
}


//checkIfExist return true or false depending if the comparator 1 passed as parameter exist in the db 
func checkIfExist(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) bool {
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}
	// defer stmt.Close()
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		return false
	}
	return true
}

func main() {
	databaseOpened, _ := sql.Open("sqlite3", "dataBase/forum.db")
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)

	runServer()
}
