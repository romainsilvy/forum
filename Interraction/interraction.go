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
func SuppThread(session *sessions.Session, id_th string, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	checkUser := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	checkUserIsInThread := databaseTools.SingleRowQuerry(database, "id_user", "Thread", "id_th", id_th)
	conv_id_th, _ := strconv.Atoi(id_th)

	if checkUser == checkUserIsInThread {
		_, err := database.Exec(`DELETE FROM Thread WHERE id_th = ? `, conv_id_th)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Ce thread ne vous appartient pas")
	}
}

func ManageLike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
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

func ManageDislike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
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

// Fetch data from the js script on index.html
func FetchLike(db *sql.DB) {
	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		var myParam MyBody
		var item [2]int
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &myParam)
		sessionCookieAuth, _ := store.Get(r, "auth")
		littlecookie := sessionCookieAuth.Values["user"]
		user_name := fmt.Sprintf("%v", littlecookie)
		id_user := databaseTools.SingleRowQuerry(db, "id_user", "User", "user_name", user_name)
		id_user_int, _ := strconv.Atoi(id_user)
		id_th_int, _ := strconv.Atoi(myParam.Id_th)
		value_int, _ := strconv.Atoi(myParam.Value)

		if sessionCookieAuth.Values["authenticated"] == true {
			switch value_int {
			case 1:
				ManageLike(sessionCookieAuth, db, id_user_int, id_th_int)
				fmt.Println("apres le managelike")
			case -1:
				ManageDislike(sessionCookieAuth, db, id_user_int, id_th_int)
			}
		}

		item[0] = databaseTools.CountOfLike(db, myParam.Id_th, 1)
		item[1] = databaseTools.CountOfLike(db, myParam.Id_th, -1)
		bytes, _ := json.Marshal(item)

		w.Write(bytes)
	})
}
