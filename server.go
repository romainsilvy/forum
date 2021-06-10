package main

// idée a faire : mettre la base de donnée dans sql web et faire la commande ligne 69 ^^ INSHALLAH CA MARCHE
import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	databaseTools "test/dataBase"

	// "forum/dataBase/database"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleAccueil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")

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

		if r.FormValue("connect") != "" {
			connexionUser := r.FormValue("connexionUser")
			connexionPassword := r.FormValue("connexionPassword")
			checkIfExist := checkIfExist(database, "password", "User", "user_name", connexionUser)

			if checkIfExist {
				userPassword := singleRowQuerry(database, "password", "User", "user_name", connexionUser)
				if userPassword != "notExist" {
					if userPassword == connexionPassword {
						session, _ := store.Get(r, "auth")

						// Authentication goes here
						// ...

						// Set user as authenticated
						session.Values["authenticated"] = true
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

		variable.Execute(w, tabUser)
	})
}

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
		session, _ := store.Get(r, "auth")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		variable.Execute(w, result)
	})
}

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
