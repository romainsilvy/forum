package displayTools

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	databaseTools "tools/dataBase"
)

// Display the category and append values from sql Thread table
func DisplayCategory(inputCatChoisie string, dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveCategoryRows(db, inputCatChoisie)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err displaycategory")
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	return dataToSend
}

// Display the research from the searchBar and append values from sql Thread table
func DisplaySearchResult(inputSearchBar string, dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveSearchRows(db, inputSearchBar)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err displaysearchresult")
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	return dataToSend
}

// return a tab with all threads
func RetrieveAccueil(dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveAccueilRows(db)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err retrieveaccueil")
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	for _, content := range dataToSend {
		databaseTools.CheckIfExistLike(db, content.Id_th, content.Id_user)
		content.Like = databaseTools.CountOfLike(db, strconv.Itoa(content.Id_th), 1)
	}
	return dataToSend
}

func DisplayThread(db *sql.DB, id_th int64) databaseTools.ThreadData {
	item := databaseTools.ThreadData{}
	db.QueryRow(`SELECT * FROM Thread WHERE id_th = ?`, id_th).Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
	return item
}
