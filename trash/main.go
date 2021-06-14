package tryhard

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

<<<<<<< HEAD
// 	var result []TestItem
// 	for rows.Next() {
// 		item := TestItem{}
// 		err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email)
// 		if err2 != nil {
// 			panic(err2)
// 		}
// 		result = append(result, item)
// 	}
// 	fmt.Println(result)
// 	return result
// }

// get Data from the sqlite database and print them int the html page
// func getIncrisption(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		variable, _ := template.ParseFiles("index.html")
// 		http.SetCookie(w, &http.Cookie{
// 			Name:  "my-cookie",
// 			Value: "some value",
// 			Path:  "/",
// 		})
// 		fmt.Println("COOKIE CREATED")

// 		if r.FormValue("inscriptionPseudo") != "" {
// 			inscriptionPseudo := r.FormValue("inscriptionPseudo")
// 			inscriptionEmail := r.FormValue("inscriptionEmail")
// 			inscriptionEmailConfirm := r.FormValue("inscriptionEmailConfirm")
// 			inscriptionPassword := r.FormValue("inscriptionPassword")
// 			inscriptionPasswordConfirm := r.FormValue("inscriptionPasswordConfirm")

// 			if inscriptionEmail == inscriptionEmailConfirm && inscriptionPassword == inscriptionPasswordConfirm {
// 				databaseTools.InsertIntoUsers(inscriptionPseudo, inscriptionEmail, inscriptionPassword, "test")
// 			}
// 		}
// 		rows, _ := database.Query("select * from User")
// 		result := tabUser
// 		for rows.Next() {
// 			item := oneUser
// 			err2 := rows.Scan(&item.Id_user, &item.User_name, &item.Password, &item.Email, &item.Image)
// 			if err2 != nil {
// 				panic(err2)
// 			}
// 			result = append(result, item)
// 		}
// 		variable.Execute(w, result)
// 	})
// }

// func getConnexion(oneUser databaseTools.User, tabUser []databaseTools.User, database *sql.DB) {
// 	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
// 		variable, _ := template.ParseFiles("index.html")

// 		connexionUser := r.FormValue("connexionUser")
// 		// connexionPassword := r.FormValue("connexionPassword")

// 		rows, _ := database.Query("select * from User")
// 		result := tabUser
// 		for rows.Next() {
// 			item := oneUser
// 			err2 := rows.Scan(&item.User_name)
// 			if err2 != nil {
// 				panic(err2)
// 			}
// 			if connexionUser == item.User_name {
// 				fmt.Println("c est good")
// 			} else {
// 				fmt.Println("marche pas")
// 			}
// 			result = append(result, item)
// 		}
// 		variable.Execute(w, result)
// 	})
// }
=======
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
>>>>>>> feature/thread-ouvert
