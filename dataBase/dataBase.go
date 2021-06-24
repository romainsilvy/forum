package databaseTools

import (
	"database/sql"
	"fmt"
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
	Id_th      int
	Id_user    int
	Title      string
	Content    string
	Created_at string
	Category   string
}

type ThreadData struct {
	Id_th      int
	Id_user    int
	Title      string
	Content    string
	Created_at string
	Category   string
	Like       int
	Dislike    int
}

type ThreadToSend struct {
	Id_th    int
	Id_user  int
	Title    string
	Content  string
	Category string
	Nbr_like int
}

type ThreadMessage struct {
	Id_th_msg   int
	Id_th       int
	Id_user     int
	Created_at  string
	Msg_content string
}

type Like struct {
	Id_user     int
	Nbr_like    int
	Nbr_dislike int
	Id_th       int
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

func InsertIntoLike(id_user int, id_th int, value int, db *sql.DB) {
	fmt.Println("rentre dans insertintolike")
	_, err := db.Exec(`INSERT INTO Like (id_user, id_th, value) VALUES (?, ?, ?)`, id_user, id_th, value)
	fmt.Println("apres le db exec")
	if err != nil {
		fmt.Println("err insertintolike")
		log.Fatal(err)
	}
	fmt.Println("apres le err")
}

// func InsertIntoThreads(id_user int, title string, content string, created_at string, db *sql.DB) {
// 	_, err := db.Exec(`INSERT INTO Thread (id_user, title, content, created_at, notif, like_count, dislike_count, comment_count) VALUES (?, ?, ?, ?, false, 0, 0, 0)`, id_user, title, content, created_at)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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

func SingleRowQuerryLike(db *sql.DB, comparator1 string, comparator2 int, comparator3 string, comparator4 int) string {

	stmt, err := db.Prepare("select value from Like where " + comparator1 + " = ? and " + comparator3 + " = ?")
	if err != nil {
		log.Fatal(err)
	}
	var toReturn string
	err = stmt.QueryRow(comparator2, comparator4).Scan(&toReturn)
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

func CheckIfExistDeux(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 int) bool {
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

func RetrieveCategoryRows(db *sql.DB, inputCatChoisie string) *sql.Rows {
	reqC := `SELECT 
		id_user,
		title,
		content,
		created_at,
		category
		FROM 
		Thread
		WHERE category = ?
		ORDER BY created_at DESC`
	rows, _ := db.Query(reqC, inputCatChoisie)
	return rows
}

func RetrieveSearchRows(db *sql.DB, inputSearchBar string) *sql.Rows {
	reqS := `SELECT 
			id_user,
			title,
			content,
			created_at,
			category
			FROM 
			Thread
			WHERE title = ?
			ORDER BY created_at DESC`
	rows, _ := db.Query(reqS, inputSearchBar)
	return rows
}

func RetrieveAccueilRows(db *sql.DB) *sql.Rows {
	reqS := `SELECT 
			id_th,
			id_user,
			title,
			content,
			created_at,
			category
			FROM 
			Thread
			ORDER BY created_at DESC`
	rows, _ := db.Query(reqS)
	return rows
}

func CountOfLike(db *sql.DB, id_th string, value int) int {
	req := `SELECT
			COUNT(*)
			FROM
			Like
			Where id_th = ?
			AND 
			value = ?`
	rows := db.QueryRow(req, id_th, value)
	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0
	}

	return count

}

func CheckIfExistLike(db *sql.DB, id_th int, id_user int) bool {
	req := `SELECT
			COUNT(*)
			FROM
			Like
			Where id_th = ?
			AND
			id_user = ?`
	rows := db.QueryRow(req, id_th, id_user)
	var count int
	err := rows.Scan(&count)
	if err != nil {
		fmt.Println("checkifexistlike err")
		panic(err)
	}
	if count == 0 {
		return false
	}
	return true

}
