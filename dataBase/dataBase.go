package database

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

func InsertIntoUsers(user_name string, email string, password string, image string) {
	db := initDatabase("test.db")
	_, err := db.Exec(`INSERT INTO User (user_name, email, password, image) VALUES (?, ?, ?, ?)`, user_name, email, password, image)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// func main() {
// 	db := initDatabase("test.db")
// 	InsertIntoUsers(db, "Louis", "mail", "mdp", "image")
// 	// ReadItem(db)
// 	defer db.Close()

// }

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
