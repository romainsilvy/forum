package accountTools

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	databaseTools "tools/dataBase"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

//global vars
var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//inscription manage the inscription form
func Inscription(r *http.Request, database *sql.DB) {
	inscriptionPseudo := r.FormValue("inscriptionPseudo")

	if inscriptionPseudo != "" {
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			hashed := HashAndSalt(inscriptionPassword)
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, hashed, database)
		} else {
			fmt.Println("Le compte n'a pas pu être créé")
		}
	}
}

//connexion manage the connexion form
func Connexion(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.FormValue("connect") != "" {
		connexionUser := r.FormValue("connexionUser")
		connexionPassword := r.FormValue("connexionPassword")

		passwordHashed := databaseTools.SingleRowQuerry(database, "password", "User", "user_name", connexionUser)
		checkIfExist := databaseTools.CheckIfExist(database, "password", "User", "user_name", connexionUser)
		fmt.Println(checkIfExist)

		if checkIfExist {
			if passwordHashed != "notExist" {
				if ComparePasswords(passwordHashed, connexionPassword) {
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

func ChangePassword(r *http.Request, userPassword string, userName string, database *sql.DB) {
	changePasswordButton := r.FormValue("changePasswordButton")

	if changePasswordButton != "" {
		oldPassword := r.FormValue("oldPassword")
		if ComparePasswords(userPassword, oldPassword) {
			newPassword := r.FormValue("newPassword")
			newPasswordConfirm := r.FormValue("newPasswordConfirm")
			newPasswordHashed := HashAndSalt(newPassword)
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

func ChangeEmail(r *http.Request, userPassword string, userName string, db *sql.DB) {
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

func HashAndSalt(pwd string) string {

	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	BytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, BytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
