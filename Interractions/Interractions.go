package interractionsTools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	displayTools "tools/Display"
	databaseTools "tools/dataBase"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("ismatheplatypus@w*")
	store = sessions.NewCookieStore(key)
)

//		STRUC FOR JS LIKE / DISLIKES
type MyBody struct {
	Id_th string `json:id_th`
	Value string `json:value`
}

type contentThread struct {
	Title    string `json:title`
	Content  string `json:content`
	Category string `json:category`
}

// requete ajoute un thread
func AddThread(session *sessions.Session, title string, content string, category string, database *sql.DB) (int64, error) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	id_user_not_converted := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(id_user_not_converted)
	result, err := database.Exec(`INSERT INTO Thread (id_user, title, content,  category, created_at) VALUES (?, ?, ?, ?, time())`, id_user, title, content, category)

	if err != nil {
		log.Fatal(err)
	}
	return result.LastInsertId()
}

// requete supprimer un thread
func SuppThread(session *sessions.Session, database *sql.DB) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	checkUser := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	checkThread := databaseTools.SingleRowQuerry(database, "id_th", "Thread", "id_user", checkUser)
	id_thread, _ := strconv.Atoi(checkThread)
	_, err := database.Exec(`DELETE FROM Thread WHERE id_th = ? `, id_thread)

	if err != nil {
		log.Fatal(err)
	}
}

// add / delete / update likes in the bdd
func ManageLike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
	if databaseTools.CheckIfExistLike(db, id_th_int, id_user_int) {
		if databaseTools.SingleRowQuerryLike(db, "id_th", id_th_int, "id_user", id_user_int) == "1" {
			db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
		} else {
			db.Exec(`UPDATE Like SET value = ? WHERE id_user = ? and id_th = ?`, 1, id_user_int, id_th_int)
		}
	} else {
		databaseTools.InsertIntoLike(id_user_int, id_th_int, 1, db)
	}
}

// add / delete / update dislikes in the bdd
func ManageDislike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
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

func FetchThread(db *sql.DB) {
	http.HandleFunc("/thread", func(w http.ResponseWriter, r *http.Request) {
		var threadBody contentThread
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &threadBody)
		sessionCookieAuth, _ := store.Get(r, "auth")
		title := threadBody.Title
		content := threadBody.Content
		category := threadBody.Category

		id_th, _ := AddThread(sessionCookieAuth, title, content, category, db)
		item := displayTools.RetrieveThread(db, id_th)
		bytes, _ := json.Marshal(item)

		w.Write(bytes)
	})
}
