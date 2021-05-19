// package tryhard

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type TestItem struct {
	Id_user   string
	User_name string
	Password  string
	Email     string
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		Id TEXT NOT NULL PRIMARY KEY,
		Name TEXT,
		Phone TEXT,
		InsertedDatetime DATETIME
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func StoreItem(db *sql.DB, items []TestItem) {
	sql_additem := `
	INSERT OR REPLACE INTO items(
		Id_user,
		User_name,
		Password,
		Email,
		InsertedDatetime
	) values(?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Id_user, item.User_name, item.Password, item.Email)
		if err2 != nil {
			panic(err2)
		}
	}
}

func ReadItem(db *sql.DB) []TestItem {
	sql_readall := `
	SELECT Id_user, User_name, Password, Email FROM items
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []TestItem
	for rows.Next() {
		item := TestItem{}
		err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	fmt.Println(result)
	return result
}
