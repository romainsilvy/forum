package main 

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initDatabase(database string) *sql.DB {
    db, err := sql.Open("sqlite3", database) 
    if err != nil {
        log.Fatal(err) 
    } 

    sqltStmt := `
                CREATE TABLE IF NOT EXISTS users (
                    id INSERT TO NOT NULL PRIMARY KEY AUTOINCREMENT,
					email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL,
                )
                `
    return db 
			}
			
func insertIntoUsers(db *sql.DB, name string, email string, password string) (int64, error) {
	
	result, _ := db.Exec("INSERT INTO users (name) (email) (password) VALUES (?) (?) (?) ", name, email, password)

	return result.lastInsertId()

}