package displayTools

import (
	"database/sql"
	"html/template"
	"net/http"

	databaseTools "tools/dataBase"
)

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
	variable.Execute(w, dataToSend)
}

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
	variable.Execute(w, dataToSend)
}

func DisplayAccueil(dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveAccueilRows(db)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	variable.Execute(w, dataToSend)
}
