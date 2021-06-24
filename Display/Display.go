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
func DisplayAccueil(dataToSend []databaseTools.ThreadData, w http.ResponseWriter, db *sql.DB) []databaseTools.ThreadData {
	rows := databaseTools.RetrieveAccueilRows(db)
	fmt.Println(rows)
	fmt.Println("dans displayaccueil")

	//checkifexist avec idth
	//si existe alors on display normal
	//si existe pas alors on display 0

	//on recup tt les thread
	//on teste si like existe
	//si il existe pas alors on lui donne 0
	//si il existe on recup le nombre avec getnuberof

	// , &item.Id_user, &item.Title, &item.Content, &item.Category

	for rows.Next() {
		fmt.Println("execution de la req c'est bon")
		item := databaseTools.ThreadData{}
		err2 := rows.Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
		if err2 != nil {
			panic(err2)
		}
		dataToSend = append(dataToSend, item)
	}
	// fmt.Println("la fonction return :", dataToSend)

	for _, content := range dataToSend {
		fmt.Println(content)
		//checkifexistlike sur content
		//si le like existe
		//on le append
		//si il existe pas
		//on append like = 0
	}

	return dataToSend
}

func DisplayThread(db *sql.DB, id_th int64) databaseTools.ThreadData {
	item := databaseTools.ThreadData{}
	db.QueryRow(`SELECT * FROM Thread WHERE id_th = ?`, id_th).Scan(&item.Id_th, &item.Id_user, &item.Title, &item.Content, &item.Created_at, &item.Category)
	return item
}
