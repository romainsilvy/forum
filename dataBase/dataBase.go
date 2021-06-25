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
	OneThread []Thread
}

type Thread struct {
	Id_th   int
	Title   string
	Content string
	Like    int
	Dislike int
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

//InitDatabase is the function which init the database
func InitDatabase(database string) *sql.DB {
	//open the database in sqlite
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	//read the instructions from databaseText.txt
	createDb, err := ioutil.ReadFile("dataBase/databaseText.txt")
	if err != nil {
		log.Fatal(err)
	}

	//create the database using the instructions contained in createDb
	str := string(createDb)
	_, err = db.Exec(str)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//InsertIntoUsers is the function which insert data in the table User
func InsertIntoUsers(user_name string, email string, password string, db *sql.DB) {
	_, err := db.Exec(`INSERT INTO User (user_name, email, password) VALUES (?, ?, ?)`, user_name, email, password)
	if err != nil {
		log.Fatal(err)
	}
}

//InsertIntoLike is the function which insert data in the table Like
func InsertIntoLike(id_user int, id_th int, value int, db *sql.DB) {
	_, err := db.Exec(`INSERT INTO Like (id_user, id_th, value) VALUES (?, ?, ?)`, id_user, id_th, value)
	if err != nil {
		log.Fatal(err)
	}
}

//SingleRowQuerry is the function which select data in the table given as argument and return the wanted value
func SingleRowQuerry(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) string {
	//prepare the queryRow request
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}

	//execute the queryRow request
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		return "notExist"
	}

	return toReturn
}

func SingleRowQuerryId(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) int {
	//prepare the queryRow request
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}

	//execute the queryRow request
	var toReturn int
	err = stmt.QueryRow(comparator2).Scan(&toReturn)
	if err != nil {
		fmt.Println("err in SingleRowQuerryId", err)
	}

	return toReturn
}

//SingleRowQuerryLike is the function which select data in the table Like and return the wanted value
func SingleRowQuerryLike(db *sql.DB, comparator1 string, comparator2 int, comparator3 string, comparator4 int) string {
	//prepare the queryRow request
	stmt, err := db.Prepare("select value from Like where " + comparator1 + " = ? and " + comparator3 + " = ?")
	if err != nil {
		log.Fatal(err)
	}

	//execute the queryRow request
	var toReturn string
	err = stmt.QueryRow(comparator2, comparator4).Scan(&toReturn)
	if err != nil {
		return "notExist"
	}
	return toReturn
}

//CheckIfExist return true or false depending if the comparator 1 passed as parameter exist in the db
func CheckIfExist(db *sql.DB, rowName string, tableName string, comparator1 string, comparator2 string) bool {
	//prepare the queryRow request
	stmt, err := db.Prepare("select " + rowName + " from " + tableName + " where " + comparator1 + " = ?")
	if err != nil {
		log.Fatal(err)
	}

	//execute the queryRow request
	var toReturn string
	err = stmt.QueryRow(comparator2).Scan(&toReturn)

	//test errors
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		} else {
			return false
		}
	}
	return true
}

//UpdateValue is the function which changes the value of a row from a table
func UpdateValue(db *sql.DB, tableName string, collumnName string, newValue string, comparator1 string, comparator2 string) {
	_, err := db.Exec("update " + tableName + " set " + collumnName + " = " + "\"" + newValue + "\"" + " where " + comparator1 + " = " + "\"" + comparator2 + "\"")
	if err != nil {
		log.Fatal(err)
	}
}

//RetrieveCategoryRows is the function which retrieves threads filtered by categories
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

//RetrieveSearchRows is the function which retrieves threads filtered by the value from the searchbar
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

//RetrieveSearchRows is the function which retrieves all threads
func RetrieveAccueilRows(db *sql.DB) *sql.Rows {
	req := `SELECT 
	id_th,
	id_user,
	title,
	content,
	created_at,
	category
			FROM 
			Thread
			ORDER BY id_th DESC`
	rows, _ := db.Query(req)
	return rows
}

func RetrieveThreadcreeRow(db *sql.DB, id_user int) *sql.Rows {
	req := `SELECT 
	id_th,
	title,
	content
	FROM 
	Thread
	Where id_user = ?
	ORDER BY id_th DESC`
	rows, err := db.Query(req, id_user)
	if err != nil {
		fmt.Println("err in retrievethreadcreerow")
	}
	return rows
}

//CountOfLike is the function which returns the number of likes of a thread
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
		fmt.Println("pas de like")
		return 0
	}
	return count
}

//CheckIfExistLike is the function which returns true or false depending of the number of like (0 is false and > 0 is true)
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
		fmt.Println("err in CheckIfExistLike : ", err)
	}
	if count == 0 {
		return false
	}
	return true
}

//CheckIfThread is the function which returns true or false depending of the number of like (0 is false and > 0 is true)
func CheckIfThread(db *sql.DB, id_user int) bool {

	req := `SELECT
			id_th 
			FROM
			Thread
			Where id_user = ?`
	rows, _ := db.Query(req)

	if rows != nil {
		fmt.Println("Pas de thread")
		return false
	}
	return true
}
