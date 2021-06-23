package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	databaseTools "test/dataBase"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//global vars
var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//this structure is used to unmarshall the value of id_th
type MyBody struct {
	Id_th string `json:id_th`
	Value string `json:value`
}

//inscription manage the inscription form
func inscription(r *http.Request, database *sql.DB) {
	inscriptionPseudo := r.FormValue("inscriptionPseudo")

	if inscriptionPseudo != "" {
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			hashed := hashAndSalt(inscriptionPassword)
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, hashed, database)
		} else {
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

		//session cookie
		sessionCookieAuth, _ := store.Get(r, "auth")

		inputSearchBar := r.FormValue("searchWord")
		inputCatThread := r.FormValue("drone")

		if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] == true) {
			addThread(sessionCookieAuth, title, content, inputCatThread, database)
		} else if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] != true) {
			fmt.Println("Veuillez vous connecter pour poster un thread !")
		}

		connexion(w, r, database)
		inscription(r, database)

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

// requete ajoute un tread
func addThread(session *sessions.Session, title string, content string, category string, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	check := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(check)
	_, err := database.Exec(`INSERT INTO Thread (id_user, title, content,  category, created_at) VALUES (?, ?, ?, ?, time())`, id_user, title, content, category)
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
func handleProfil(oneUser databaseTools.User, database *sql.DB) {
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
	handleAccueil(db)
	handleProfil(databaseTools.User{}, db)
	FetchLike(db)

}

func manageLike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int, value_int int) {
	if sessionCookieAuth.Values["authenticated"] == true {
		if databaseTools.CheckIfExistLike(db, "id_user", id_user_int) && databaseTools.CheckIfExistLike(db, "id_th", id_th_int) && databaseTools.CheckIfExistLike(db, "value", value_int) {
			db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
		} else {
			databaseTools.InsertIntoLike(id_user_int, id_th_int, value_int, db)
		}
	} else if sessionCookieAuth.Values["authenticated"] != true {
		fmt.Println("Veuillez vous connecter pour poster un thread !")
	}
}

func manageDisike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int, value_int int) {
	if sessionCookieAuth.Values["authenticated"] == true {
		if databaseTools.CheckIfExistLike(db, "id_user", id_user_int) && databaseTools.CheckIfExistLike(db, "id_th", id_th_int) && databaseTools.CheckIfExistLike(db, "value", value_int) {
			db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
		} else {
			databaseTools.InsertIntoLike(id_user_int, id_th_int, value_int, db)
		}
	} else if sessionCookieAuth.Values["authenticated"] != true {
		fmt.Println("Veuillez vous connecter pour poster un thread !")
	}
}

func FetchLike(db *sql.DB) {
	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		//insere un like en fonction du post id

		var myParam MyBody

		body, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(body, &myParam)

		sessionCookieAuth, _ := store.Get(r, "auth")
		littlecookie := sessionCookieAuth.Values["user"]
		user_name := fmt.Sprintf("%v", littlecookie)
		id_user := databaseTools.SingleRowQuerry(db, "id_user", "User", "user_name", user_name)
		id_user_int, _ := strconv.Atoi(id_user)
		id_th_int, _ := strconv.Atoi(myParam.Id_th)
		value_int, _ := strconv.Atoi(myParam.Value)

		switch value_int {
		case 1:
			manageLike(sessionCookieAuth, db, id_user_int, id_th_int, value_int)
		case -1:
			manageDisike(sessionCookieAuth, db, id_user_int, id_th_int, value_int)
		}

		databaseTools.SendNumberOfLike(db, myParam.Id_th, w, value_int)
	})
	// recup la donner envoyer en js pour le mettre dans la base de données
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
	runServer()
}
