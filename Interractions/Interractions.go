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

// requete ajoute un thread
func AddThread(session *sessions.Session, title string, content string, category string, database *sql.DB) (int64, error) {
	littlecookie := session.Values["user"]
	convertissor := fmt.Sprintf("%v", littlecookie)
	check := databaseTools.SingleRowQuerry(database, "id_user", "User", "user_name", convertissor)
	id_user, _ := strconv.Atoi(check)
	res, err := database.Exec(`INSERT INTO Thread (id_user, title, content,  category, created_at) VALUES (?, ?, ?, ?, time())`, id_user, title, content, category)
	if err != nil {
		log.Fatal(err)
	}
	return res.LastInsertId()
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

// Allow user to add / delete likes
func ManageLike(sessionCookieAuth *sessions.Session, db *sql.DB, id_user_int int, id_th_int int) {
	if sessionCookieAuth.Values["authenticated"] == true {
		fmt.Println("user connecté")
		if databaseTools.CheckIfExistLike(db, id_th_int, id_user_int) {
			fmt.Println("like existe")
			if databaseTools.SingleRowQuerryLike(db, "id_th", id_th_int, "id_user", id_user_int) == "1" {
				fmt.Println("c'est un like")
				db.Exec(`DELETE FROM Like WHERE id_user = ? AND id_th = ?`, id_user_int, id_th_int)
			} else {
				fmt.Println("c'est un dislike ")
				db.Exec(`UPDATE Like SET value = ? WHERE id_user = ? and id_th = ?`, 1, id_user_int, id_th_int)
			}
		} else {
			fmt.Println("like existe pas ")
			databaseTools.InsertIntoLike(id_user_int, id_th_int, 1, db)
			fmt.Println("apres le insertintolike")
		}
	}
}

// Allow user to add / delete dislikes
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
			ManageLike(sessionCookieAuth, db, id_user_int, id_th_int)
			fmt.Println("apres le managelike")
		case -1:
			ManageDislike(sessionCookieAuth, db, id_user_int, id_th_int)
		}

		dislike := databaseTools.CountOfLike(db, myParam.Id_th, -1)
		fmt.Println("apres count dislike")
		like := databaseTools.CountOfLike(db, myParam.Id_th, 1)
		fmt.Println("apres count like")

		var item [2]int
		item[0] = like
		item[1] = dislike
		bytes, _ := json.Marshal(item)

		w.Write(bytes)
	})
}

type contentThread struct {
	Title    string `json:title`
	Content  string `json:content`
	Category string `json:category`
}

func FetchThread(db *sql.DB) {
	fmt.Println("dans le FetchThread")
	http.HandleFunc("/thread", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("dans le HundleFunc")
		var threadBody contentThread
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &threadBody)
		sessionCookieAuth, _ := store.Get(r, "auth")
		// user_name := fmt.Sprintf("%v", littlecookie)

		// id_user := databaseTools.SingleRowQuerry(db, "id_user", "User", "user_name", user_name)
		title := threadBody.Title
		content := threadBody.Content
		category := threadBody.Category

		id_th, _ := AddThread(sessionCookieAuth, title, content, category, db) /* ---FONCTIONNEL */
		item := displayTools.DisplayThread(db, id_th)
		bytes, _ := json.Marshal(item)
		w.Write(bytes)
	})
}

// au début charge les thread avec le template								DONE
// pour ajouter thread on fait un last insert id en sql
// on balance les data dans la route /thread
// on fetch en js
// dans le response on add en appenchild dans la partie div des threads
// puis on affiche sur la page
