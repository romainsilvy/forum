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

func handleAccueil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("index.html")
		http.SetCookie(w, &http.Cookie{
			Name:  "my-cookie",
			Value: "some value",
			Path:  "/",
		})

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
						fmt.Println("tu est co chacal")
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

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

// get Data from the sqlite database and print them int the html page
// func getIncrisption(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		variable, _ := template.ParseFiles("index.html")
// 		http.SetCookie(w, &http.Cookie{
// 			Name:  "my-cookie",
// 			Value: "some value",
// 			Path:  "/",
// 		})
// 		fmt.Println("COOKIE CREATED")

// 		if r.FormValue("inscriptionPseudo") != "" {
// 			inscriptionPseudo := r.FormValue("inscriptionPseudo")
// 			inscriptionEmail := r.FormValue("inscriptionEmail")
// 			inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
// 			inscriptionPassword := r.FormValue("inscriptionPassword")
// 			inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

// 			if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
// 				databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "test")
// 			}
// 		}
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

// func getConnexion(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
// 	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
// 		variable, _ := template.ParseFiles("index.html")

// 		connexionUser := r.FormValue("connexionUser")
// 		// connexionPassword := r.FormValue("connexionPassword")

// 		rows, _ := database.Query("select * from User")
// 		result := tabUser
// 		for rows.Next() {
// 			item := oneUser
// 			err2 := rows.Scan(&item.User_name)
// 			if err2 != nil {
// 				panic(err2)
// 			}
// 			if connexionUser == item.User_name {
// 				fmt.Println("c est good")
// 			} else {
// 				fmt.Println("marche pas")
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
	// handleProfil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
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
