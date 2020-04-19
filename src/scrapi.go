package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host     = "goscraper-db.cnwfupzwzfus.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "scrapi"
	password = "354scrapipw"
	dbname   = "goscraper"
)

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Homepage Endpoint Hit")

}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/article", AllArticles)
	myRouter.HandleFunc("/article", AllArticles)
	fmt.Println("scraper-api started. Listening on port:8081")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func startApi(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting scraper-api")
	handleRequests()
}

