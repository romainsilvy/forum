package displayTools

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	databaseTools "tools/dataBase"
)

// return a tab with all threads
func RetrieveAccueil(dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveAccueilRows(db)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err in RetrieveAccueil : ", err2)
		}
		dataToSend = append(dataToSend, item)
	}
	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}
	return dataToSend
}

func DisplayAccueil(dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveAccueilRows(db)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err in DisplayAccueil : ", err2)
		}
		dataToSend = append(dataToSend, item)
	}
	variable.Execute(w, dataToSend)
}

func DisplayThreadCree(Threadcree []databaseTools.Thread, db *sql.DB, id_user int) []databaseTools.Thread {
	rows := databaseTools.RetrieveThreadcreeRow(db, id_user)
	for rows.Next() {
		item := databaseTools.Thread{}
		err2 := rows.Scan(&item.Id_th, &item.Title, &item.Content)
		if err2 != nil {
			fmt.Println("err in DisplayThreadCree : ", err2)
		}
		Threadcree = append(Threadcree, item)
	}
	// for i := 0; i < len(Threadcree); i++ {
	// 	Threadcree[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(Threadcree[i].Id_th), 1)
	// 	Threadcree[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(Threadcree[i].Id_th), -1)
	// }
	return Threadcree
}

func DisplayCategory(inputCatChoisie string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveCategoryRows(db, inputCatChoisie)

	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println()
		}
		dataToSend = append(dataToSend, item)
	}
	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}
	variable.Execute(w, dataToSend)
}

func DisplaySearchResult(inputSearchBar string, dataToSend []databaseTools.ThreadData, variable *template.Template, w http.ResponseWriter, db *sql.DB) {
	rows := databaseTools.RetrieveSearchRows(db, inputSearchBar)
	for rows.Next() {
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			fmt.Println("err in DisplaySearchResult : ", err2)
		}
		dataToSend = append(dataToSend, item)
	}
	for i := 0; i < len(dataToSend); i++ {
		dataToSend[i].Like = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), 1)
		dataToSend[i].Dislike = databaseTools.CountOfLike(db, strconv.Itoa(dataToSend[i].Id_th), -1)
	}
	variable.Execute(w, dataToSend)
}
