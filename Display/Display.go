package displayTools

import (
	"database/sql"
	"fmt"
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
			fmt.Println("err retrieveaccueil")
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

func RetrieveThread(db *sql.DB, id_th int64) databaseTools.ThreadData {
	item := databaseTools.ThreadData{}
	db.QueryRow(`SELECT * FROM Thread WHERE id_th = ?`, id_th).Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
	return item
}
