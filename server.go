package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	databaseTools "test2/dataBase"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//inscription manage the inscription form
func inscription(r *http.Request) {
	inscriptionPseudo := r.FormValue("inscriptionPseudo")
	inscriptionEmail := r.FormValue("inscriptionEmail")
	inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
	inscriptionPassword := r.FormValue("inscriptionPassword")
	inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")
	if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
		hashed := hashAndSalt(inscriptionPassword)
		databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, hashed)
	}
}

//connexion manage the connexion form
func connexion(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.FormValue("connect") != "" {
		connexionUser := r.FormValue("connexionUser")
		connexionPassword := r.FormValue("connexionPassword")

		passwordHashed := databaseTools.SingleRowQuerry(database, "password", "User", "user_name", connexionUser)
		checkIfExist := databaseTools.CheckIfExist(database, "password", "User", "user_name", connexionUser)

		if checkIfExist {
			if passwordHashed != "notExist" {
				if comparePasswords(passwordHashed, connexionPassword) {
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

func changePassword(r *http.Request, userPassword string, userName string, database *sql.DB) {
	changePasswordButton := r.FormValue("changePasswordButton")

	if changePasswordButton != "" {
		oldPassword := r.FormValue("oldPassword")
		if comparePasswords(userPassword, oldPassword) {
			newPassword := r.FormValue("newPassword")
			newPasswordConfirm := r.FormValue("newPasswordConfirm")
			newPasswordHashed := hashAndSalt(newPassword)
			// newPasswordConfirmHashed := hashAndSalt(newPasswordConfirm)

			if newPassword == newPasswordConfirm {
				databaseTools.UpdateValue(database, "User", "password", newPasswordHashed, "user_name", userName)
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

func hashAndSalt(pwd string) string {

	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	BytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, BytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(db)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, db)
}

//------------------------------------------------- A REFAIRE -------------------------------------------------
//handleAccueil is the handlefunc for the main page
func handleAccueil(database *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var dataToSend databaseTools.Data
		variable, _ := template.ParseFiles("index.html")
		title := r.FormValue("threadTitle")
		thread := r.FormValue("créa_thread")
		addThread(r, databaseTools.User{}, title, thread, database)
		inscriptionEmail := r.FormValue("inscriptionEmail")
		if inscriptionEmail == "" {
			return
		} else {
			fmt.Println("wtf")
			inscription(r)
		}
		connexion(w, r, database)

		// req := `SELECT Thread.id_user,
		// Thread.title,
		// Thread.content,
		// User.user_name
		// FROM Thread, User`
		// rows, _ := database.Query(req)
		// fmt.Println(rows)
		// for rows.Next() {
		// 	item := databaseTools.ThreadData{}
		// 	fmt.Println("dans le row next")
		// 	fmt.Println(item)
		// 	err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.User_name)
		// 	if err2 != nil {
		// 		panic(err2)
		// 	}
		// 	dataToSend.Posts = append(dataToSend.Posts, item)
		// }

		req := `SELECT Thread.Content FROM Thread`
		rows, _ := database.Query(req)
		for rows.Next() {
			item := databaseTools.ThreadData{}
			err2 := rows.Scan(&item.Content)
			if err2 != nil {
				panic(err2)
			}
			fmt.Println(item)
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")

			dataToSend.Posts = append(dataToSend.Posts, item)
		}

		fmt.Println(dataToSend.Posts)
		variable.Execute(w, dataToSend)
	})
}

func addThread(r *http.Request, oneUser databaseTools.User, title string, content string, database *sql.DB) {
	if r.FormValue("créa_thread") == "Submit" {
		idUser := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", oneUser.User_name)
		id, _ := strconv.Atoi(idUser)
		databaseTools.InsertIntoThreads(id, title, content, "10/06/21 10:35", database)
	}
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	databaseTools.InitDatabase("dataBase/forum.db")
	db, _ := sql.Open("sqlite3", "dataBase/forum.db")
	handleAll(db)
	runServer()
}
