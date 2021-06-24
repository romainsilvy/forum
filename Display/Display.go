package displayTools

import (
	"database/sql"
	"fmt"
	"net/http"
	databaseTools "tools/dataBase"
)

// Display the category and append values from sql Thread table
func DisplayCategory(inputCatChoisie string, dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveCategoryRows(db, inputCatChoisie)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
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
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	return dataToSend
}

// Display all the threads on the home page and append the values in sql Thread Table
func DisplayAccueil(dataToSend []databaseTools.ThreadToSend, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadToSend {
	rows := databaseTools.RetrieveAccueilRows(db)
	fmt.Println("creation de la req c'est bon")
	for rows.Next() {
		fmt.Println("execution de la req c'est bon")
		item := databaseTools.ThreadToSend{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Category, &item.Nbr_like)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	fmt.Println("la fonction return :", dataToSend)
	return dataToSend
}

func DisplayThread(db *sql.DB, id_th int64) databaseTools.ThreadData {
	item := databaseTools.ThreadData{}
	db.QueryRow(`SELECT * FROM Thread WHERE id_th = ?`, id_th).Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
	return item
}
