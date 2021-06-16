package databaseTools

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id_user   int
	User_name string
	Password  string
	Email     string
	Image     string
}

type Thread struct {
	Id_th         int
	Id_user       int
	Title         string
	Content       string
	Created_at    string
	Notif         bool
	Like_count    int
	Dislike_count int
	Comment_count int
}

type ThreadData struct {
	Id_user int
	Title   string
	Content string
}

type ThreadMessage struct {
	Id_th_msg   int
	Id_th       int
	Id_user     int
	Created_at  string
	Msg_content string
}

type Like struct {
	Id_like int
	Id_user int
	Like    bool
	Dislike bool
	Id_th   int
}

type Following struct {
	Id_following int
	Id_user      int
	Id_th        int
}

type ThreadCategory struct {
	Id_th_cat int
	Id_cat    int
	Id_th     int
}

type Category struct {
	Id_cat   int
	Cat_name string
}

func InitDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	createDb, err := ioutil.ReadFile("dataBase/databaseText.txt")
	if err != nil {
		log.Fatal(err)
	}
	str := string(createDb)
	_, err = db.Exec(str)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InsertIntoUsers(user_name string, email string, password string, db *sql.DB) {
	_, err := db.Exec(`INSERT INTO User (user_name, email, password) VALUES (?, ?, ?)`, user_name, email, password)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertIntoThreads(id_user int, title string, content string, created_at string, db *sql.DB) {
	_, err := db.Exec(`INSERT INTO Thread (id_user, title, content, created_at, notif, like_count, dislike_count, comment_count) VALUES (?, ?, ?, ?, false, 0, 0, 0)`, id_user, title, content, created_at)
	if err != nil {
		log.Fatal(err)
	}
}

//singleRowQuerry retrieve a value in the db with a where comparator
func SingleRowQuerry(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) string {
	// SELECT password FROM User WHERE User_name = " ?"
	//recup le mdp
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		return "notExist"
	}
	return toReturn
}

//checkIfExist return true or false depending if the comparator 1 passed as parameter exist in the db
func CheckIfExist(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) bool {
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		} else {
			return false
		}
	}
	return true
}

//UpdateValue change the value of a case
func UpdateValue(db *sql.DB, tableName string, collumnName string, newValue string, comparator1 string, comparator2 string) {
	_, err := db.Exec("update " + tableName + " set " + collumnName + " = " + "\"" + newValue + "\"" + " where " + comparator1 + " = " + "\"" + comparator2 + "\"")
	if err != nil {
		log.Fatal(err)
	}
}
