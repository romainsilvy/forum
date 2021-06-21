package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	databaseTools "test/dataBase"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//inscription manage the inscription form
func inscription(r *http.Request, database *sql.DB) {
	goodCreation := false
	inscriptionPseudo := r.FormValue("inscriptionPseudo")
	if inscriptionPseudo != "" {
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			hashed := hashAndSalt(inscriptionPassword)
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, hashed, database)
			goodCreation = true
		}

		if !goodCreation {
			fmt.Println("Le compte n'a pas pu être créé")
		}
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
					fmt.Println(session.Values["user"])
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
		var dataToSend []databaseTools.ThreadData
		variable, _ := template.ParseFiles("index.html")
		title := r.FormValue("threadTitle")
		content := r.FormValue("créa_thread")
		sub := r.FormValue("submitthread")
		inputBar := r.FormValue("searchWord")
		inputCatThread := r.FormValue("drone")
		inputCatCham := r.FormValue("CHAMEAU")
		inputCatDrom := r.FormValue("DROMADAIRE")
		inputCatLama := r.FormValue("LAMA")
		inputCatChoisie := ""
		run := false
		fmt.Println(inputCatThread)
		// catCham := ""
		session, _ := store.Get(r, "auth")
		if (sub == "Enregistrer") && (session.Values["authenticated"] == true) {
			addThread(session, title, content, inputCatThread, database)
		} else if (sub == "Enregistrer") && (session.Values["authenticated"] != true) {
			fmt.Println("Veuillez vous connecter pour poster un thread !")
		}
		connexion(w, r, database)
		inscription(r, database)

		if inputCatCham != "" {
			inputCatChoisie = inputCatCham
			run = true
		} else if inputCatDrom != "" {
			inputCatChoisie = inputCatDrom
			run = true
		} else if inputCatLama != "" {
			inputCatChoisie = inputCatLama
			run = true
		}
		fmt.Println(inputCatChoisie)
		if run {
			reqC := `SELECT 
			id_user,
			title,
			content,
			created_at,
			category
			FROM 
			Thread
			WHERE category = ?
			ORDER BY created_at DESC`
			rows, _ := database.Query(reqC, inputCatChoisie)
			for rows.Next() {
				item := databaseTools.ThreadData{}
				err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
				if err2 != nil {
					panic(err2)
				}
				fmt.Println("jsuis rentré dedans")
				dataToSend = append(dataToSend, item)
			}
			fmt.Println("good")
			variable.Execute(w, dataToSend)
		} else if inputBar != "" {
			reqS := `SELECT 
			id_user,
			title,
			content,
			created_at,
			category
			FROM 
			Thread
			WHERE title = ?
			ORDER BY created_at DESC`
			rows, _ := database.Query(reqS, inputBar)
			for rows.Next() {
				item := databaseTools.ThreadData{}
				err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
				if err2 != nil {
					panic(err2)
				}
				dataToSend = append(dataToSend, item)
			}
			variable.Execute(w, dataToSend)
		} else {
			req := `SELECT 
			id_user,
			title,
			content,
			created_at,
			category
			FROM 
			Thread
			ORDER BY created_at DESC`
			rows, _ := database.Query(req)
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
	})
}

// requete ajoute un tread
func addThread(session *sessions.Session, title string, content string, category string, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	check := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(check)
	_, err := database.Exec(`INSERT INTO Thread (id_user, title, content,  category, created_at, notif, like_count, dislike_count, comment_count) VALUES (?, ?, ?, ?, time(), false, 0, 0, 0)`, id_user, title, content, category)
	if err != nil {
		log.Fatal(err)
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

func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleAccueil(databaseTools.User{}, []databaseTools.User{}, db)
	handleProfil(databaseTools.User{}, []databaseTools.User{}, db)
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
	// databaseTools.InsertIntoThreads(10, "mon histoire", "blablabla", "crée le blabla", db)
	runServer()
}
