package displayTools

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	databaseTools "tools/dataBase"
)

//RetrieveAccueil is the function which retrieves the data to be displayed on the main page
func RetrieveAccueil(dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveAccueilRows(db)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}

	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}

	return dataToSend
}

//DisplayCategory is the function which displays the data related with categories
func DisplayCategory(inputCatChoisie string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveCategoryRows(db, inputCatChoisie)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}

	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}

	variable.Execute(w, dataToSend)
}

//DisplaySearchResult is the function which displays the data related with the searchbar
func DisplaySearchResult(inputSearchBar string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveSearchRows(db, inputSearchBar)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}

	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}

	variable.Execute(w, dataToSend)
}
