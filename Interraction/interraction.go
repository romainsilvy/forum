package interractionTools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	databaseTools "tools/dataBase"

	"github.com/gorilla/sessions"
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

// requete ajoute un thread
func AddThread(session *sessions.Session, title string, content string, category string, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	check := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(check)
	_, err := database.Exec(`INSERT INTO Thread (id_user, title, content,  category, created_at) VALUES (?, ?, ?, ?, time())`, id_user, title, content, category)
	if err != nil {
		log.Fatal(err)
	}
}

// requete supprimer un thread
func SuppThread(session *sessions.Session, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	checkUser := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(checkUser)

	checkThread := databaseTools.SingleRowQuerry(database, "id_th", "Thread", "id_user", checkUser)
	id_thread, _ := strconv.Atoi(checkThread)
	fmt.Println(id_user)
	fmt.Println(id_thread)
	_, err := database.Exec(`DELETE FROM Thread WHERE id_th = ? `, id_thread)
	if err != nil {
		log.Fatal(err)
	}
}

func manageLike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
	if sessionCookieAuth.Values["authenticated"] == true {

		if databaseTools.CheckIfExistLike(db, id_th_int, id_user_int) {
			if databaseTools.SingleRowQuerryLike(db, "id_th", id_th_int, "id_user", id_user_int) == "1" {
				db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
				fmt.Println("remove")
			} else {
				db.Exec(`UPDATE Like SET value = ? WHERE id_user = ? and id_th = ?`, 1, id_user_int, id_th_int)
				fmt.Println("modif")
			}
		} else {
			databaseTools.InsertIntoLike(id_user_int, id_th_int, 1, db)
			fmt.Println("like cree")
		}
	}
}

func manageDislike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
	if sessionCookieAuth.Values["authenticated"] == true {
		if databaseTools.CheckIfExistLike(db, id_th_int, id_user_int) {
			if databaseTools.SingleRowQuerryLike(db, "id_th", id_th_int, "id_user", id_user_int) == "-1" {
				db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
			} else {
				db.Exec(`UPDATE Like SET value = ? WHERE id_user = ? and id_th = ?`, -1, id_user_int, id_th_int)
			}
		} else {
			databaseTools.InsertIntoLike(id_user_int, id_th_int, -1, db)
		}
	}
}

func FetchLike(db *sql.DB) {
	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {

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
			manageLike(sessionCookieAuth, db, id_user_int, id_th_int)
		case -1:
			manageDislike(sessionCookieAuth, db, id_user_int, id_th_int)
		}

		dislike := databaseTools.CountOfLike(db, myParam.Id_th, -1)
		like := databaseTools.CountOfLike(db, myParam.Id_th, 1)

		w.Write([]byte(like + ":" + dislike))
	})
}
