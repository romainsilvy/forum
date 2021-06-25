package handleTools

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	accountTools "tools/AccountManagement"
	displayTools "tools/Display"
	interractionTools "tools/Interraction"
	databaseTools "tools/dataBase"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//handleAccueil is the handlefunc for the main page
func HandleAccueil(database *sql.DB) {
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

		//supp variable
		deleteButton := r.FormValue("suppr")
		fmt.Println(deleteButton)

		//session cookie
		sessionCookieAuth, _ := store.Get(r, "auth")

		inputSearchBar := r.FormValue("searchWord")
		inputCatThread := r.FormValue("drone")

		if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] == true) {
			interractionTools.AddThread(sessionCookieAuth, title, content, inputCatThread, database)
		} else if (submitButton == "Enregistrer") && (sessionCookieAuth.Values["authenticated"] != true) {
			fmt.Println("Veuillez vous connecter pour poster un thread !")
		}

		if (deleteButton != "") && (sessionCookieAuth.Values["authenticated"] == true) {
			fmt.Println("l id du post est " + deleteButton)
			interractionTools.SuppThread(sessionCookieAuth, deleteButton, database)
		}

		accountTools.Connexion(w, r, database)
		accountTools.Inscription(r, database)

		if inputCatCham != "" {
			displayTools.DisplayCategory(inputCatCham, dataToSend, variable, w, database)
		} else if inputCatDrom != "" {
			displayTools.DisplayCategory(inputCatDrom, dataToSend, variable, w, database)
		} else if inputCatLama != "" {
			displayTools.DisplayCategory(inputCatLama, dataToSend, variable, w, database)
		} else if inputSearchBar != "" {
			displayTools.DisplaySearchResult(inputSearchBar, dataToSend, variable, w, database)
		} else {
			dataToSend = displayTools.RetrieveAccueil(dataToSend, w, database)
			variable.Execute(w, dataToSend)
		}
	})
}

//handleProfil is the handlefunc for the profil page
func HandleProfil(oneUser databaseTools.User, database *sql.DB) {
	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		variable, _ := template.ParseFiles("profil.html")

		session, _ := store.Get(r, "auth")
		username := session.Values["user"].(string)
		oneUser.User_name = username
		oneUser.Email = databaseTools.SingleRowQuerry(database, "email", "User", "user_name", oneUser.User_name)
		oneUser.Password = databaseTools.SingleRowQuerry(database, "password", "User", "user_name", oneUser.User_name)
		oneUser.Id_user = databaseTools.SingleRowQuerryId(database, "id_user", "User", "user_name", oneUser.User_name)

		accountTools.ChangePassword(r, oneUser.Password, oneUser.User_name, database)
		accountTools.ChangeEmail(r, oneUser.Password, oneUser.User_name, database)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		/* JE COMMENCE A MODIFIER ICI */
		if databaseTools.CheckIfThread(database, oneUser.Id_user) { /* le check marche */
			oneUser.OneThread = displayTools.DisplayThreadCrée(oneUser.OneThread, database, oneUser.Id_user) /* ici ne marche pas */
			fmt.Println(oneUser.OneThread)
		}

		variable.Execute(w, oneUser)
	})
}
