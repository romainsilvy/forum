package main

import (
	"fmt"
	"log"
	"net/http"
)

//serveFile makes files available for the website
func serveFile() {
	// fileServer := http.FileServer(http.Dir("./data"))
}

func handle() {
	http.HandleFunc("/groupie-tracker/artist/", func(w http.ResponseWriter, r *http.Request) {

	})
}

//runServer sets the listenandserve port to 8080
func runServer() {
	fmt.Println("server is runing")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serveFile()
	runServer()
}
