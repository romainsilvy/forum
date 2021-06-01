package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id_user   int
	user_name string
	password  string
	email     string
	image     string
}

type Thread struct {
	id_th         int
	id_user       int
	title         string
	content       string
	created_at    string
	notif         bool
	like_count    int
	dislike_count int
	comment_count int
}

type ThreadMessage struct {
	id_th_msg   int
	id_th       int
	id_user     int
	created_at  string
	msg_content string
}

type Like struct {
	id_like int
	id_user int
	like    bool
	dislike bool
	id_th   int
}

type Following struct {
	id_following int
	id_user      int
	id_th        int
}

type ThreadCategory struct {
	id_th_cat int
	id_cat    int
	id_th     int
}

type Category struct {
	id_cat   int
	cat_name string
}

func initDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	createDb, _ := ioutil.ReadFile("databaseText.txt")
	str := string(createDb)
	// fmt.Println(str)
	_, err = db.Exec(str)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertIntoUsers(db *sql.DB, id_user int, user_name string, email string, password string, image string) (int64, error) {
	result, err := db.Exec(`INSERT INTO User (id_user, user_name, email, password, image) VALUES (?, ?, ?, ?, ?)`, id_user, user_name, email, password, image)
	if err != nil {
		log.Fatal(err)
	}
	return result.LastInsertId()
}

func main() {
	db := initDatabase("test.db")
	insertIntoUsers(db, 1, "Louis", "mail", "mdp", "image")
	// ReadItem(db)
	defer db.Close()

}

// func ReadItem(db *sql.DB) {
// 	sql_readall := `
// 	SELECT Id_user, User_name, Password, Email FROM items
// 	ORDER BY datetime(InsertedDatetime) DESC
// 	`

// 	rows, err := db.Query(sql_readall)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()

// 	var result []User
// 	for rows.Next() {
// 		item := User{}
// 		err2 := rows.Scan(&item.id_user, &item.user_name, &item.password, &item.email, &item.image)
// 		if err2 != nil {
// 			panic(err2)
// 		}
// 		result = append(result, item)
// 	}
// 	fmt.Println(result)
// }
