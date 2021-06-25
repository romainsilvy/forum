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

//Inscription is the function which registers a user
func Inscription(r *http.Request, database *sql.DB) {
	inscriptionPseudo := r.FormValue("inscriptionPseudo")

	//if inscriptionPseudo contains a value, it retrieves values of credentials
	if inscriptionPseudo != "" {
		inscriptionEmail := r.FormValue("inscriptionEmail")
		inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
		inscriptionPassword := r.FormValue("inscriptionPassword")
		inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

		//if the email / confirmationEmail and pasword / confirmationPasswordequals are both equal, it creates the user
		if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
			hashed := HashAndSalt(inscriptionPassword)
			databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, hashed, database)
		} else {
			fmt.Println("Le compte n'a pas pu être créé")
		}
	}
}

//Connexion is the function which registers a user
func Connexion(w http.ResponseWriter, r *http.Request, database *sql.DB) {

	//if the r.formvalue of connect contains a value, it retrieves values of credentials
	if r.FormValue("connect") != "" {
		connexionUser := r.FormValue("connexionUser")
		connexionPassword := r.FormValue("connexionPassword")
		passwordHashed := databaseTools.SingleRowQuerry(database, "password", "User", "user_name", connexionUser)

		//if there is a password associated with the user
		if databaseTools.CheckIfExist(database, "password", "User", "user_name", connexionUser) {
			//if the password is hashed
			if passwordHashed != "notExist" {
				//if password from the form and password from database are identics, it sets the session cookie
				if ComparePasswords(passwordHashed, connexionPassword) {
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

//ChangePassword is the function which update the value of the password in the database
func ChangePassword(r *http.Request, userPassword string, userName string, database *sql.DB) {
	changePasswordBtn := r.FormValue("changePasswordButton")

	//if the confirm change password button is pressed, it retrieves the value of the old password
	if changePasswordBtn != "" {
		oldPassword := r.FormValue("oldPassword")

		//if the oldpassword value and the value of the password in the database are identics, it retrieves the value of the credentials
		if ComparePasswords(userPassword, oldPassword) {
			newPassword := r.FormValue("newPassword")
			newPasswordConfirm := r.FormValue("newPasswordConfirm")
			newPasswordHashed := HashAndSalt(newPassword)

			//if both entries of the password are identics, it updates the value in the database
			if newPassword == newPasswordConfirm {
				databaseTools.UpdateValue(database, "User", "password", newPasswordHashed, "user_name", userName)
				fmt.Println("mot de passe changé ")
			} else {
				fmt.Println("confirmation de mot de passe pas bonne")
			}
		} else {
			fmt.Println("ancien mot de passe pas bon")
		}
	}
}

//ChangeEmail is the function which update the value of the email in the database
func ChangeEmail(r *http.Request, userPassword string, userName string, db *sql.DB) {
	changeEmailButton := r.FormValue("changeEmailButton")

	//if the confirm change email button is pressed, it retrieves the value of the old email
	if changeEmailButton != "" {
		newEmail := r.FormValue("newEmail")
		newEmailConfirm := r.FormValue("newEmailConfirm")

		//if the old email value and the value of the email in the database are identics, it retrieves the value of the credentials
		if newEmail == newEmailConfirm {
			password := r.FormValue("password")
			//if the password is correct, it updates the value of the email in the database
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

//HashAndSalt is the function which generate a hash of the password
func HashAndSalt(password string) string {
	passwordByte := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

//ComparePasswords is the function which compares passwords and returns true or false
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
