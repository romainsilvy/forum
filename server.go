package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	handleTools "tools/Handle"
	interractionTools "tools/Interraction"
	databaseTools "tools/dataBase"

	_ "github.com/mattn/go-sqlite3"
)

func handleAll(db *sql.DB) {
	fileServer := http.FileServer(http.Dir("./data"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handleTools.HandleAccueil(db)
	handleTools.HandleProfil(databaseTools.User{}, db)
	interractionTools.FetchLike(db)
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func main() {
	databaseTools.InitDatabase("dataBase/forum.db")
	db, _ := sql.Open("sqlite3", "dataBase/forum.db")
	handleAll(db)
	runServer()
}
