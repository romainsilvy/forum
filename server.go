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

func hanldeAccueil(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
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

		var name string
		err = db.QueryRow("SELECT id_user FROM User WHERE user_name = "Louis" ").Scan(&item.User_name)
		if err != nil {
    	log.Fatal(err)
	}
		fmt.Println(name)
		// connexionUser := r.FormValue("connexionUser")
		// connexionPassword := r.FormValue("connexionPassword")
		// seConnecter := r.FormValue("connexionUser")
		rows, _ := database.Query("select * from User")
		result := tabUser
		// aprint := tabUser
		for rows.Next() {
			item := oneUser
			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
			if err2 != nil {
				panic(err2)
			}
			// if seConnecter != "" {
			// 	if connexionUser == item.User_name && connexionPassword == item.Password {
			// 		sql_readall := `
			//  		SELECT id_user FROM User WHERE user_name = "Louis"
			//  			`
			// 		database.Query(sql_readall)
			// 		fmt.Println(database.Query(sql_readall))
			// 		row, err := database.Query(sql_readall)
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 		for row.Next() {
			// 			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
			// 			if err2 != nil {
			// 				panic(err2)
			// 			}
			// 			aprint = append(result, item)
			// 		}
			// 		fmt.Println(aprint)
			// 		// 	} else {
			// 		// 		fmt.Println("L")
			// 	}
			// }
			// result = append(result, item)
		}
		variable.Execute(w, result)
	})
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

func main() {
	databaseOpened, _ := sql.Open("sqlite3", "dataBase/forum.db")
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	hanldeAccueil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	// getIncrisption(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
	// getConnexion(databaseTools.User{}, []databaseTools.User{}, databaseOpened)
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
