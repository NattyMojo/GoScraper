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
	user	= "hamatitio"
	password = "goscrapi"
	dbname   = "goscraper"
)

var (
	Db *sql.DB
)

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Homepage Endpoint Hit")

}

func handleRequests() {
	rt := mux.NewRouter().StrictSlash(true)
	rt.HandleFunc("/api", homePage)
	rt.HandleFunc("/api/sources", GetSources).Methods("GET")
	rt.HandleFunc("/api/sources/{name}", GetSource).Methods("GET")
	rt.HandleFunc("/api/sources", CreateSource).Methods("POST")
	rt.HandleFunc("/api/articles", GetArticles).Methods("GET")
	rt.HandleFunc("/api/articles/{address}", GetArticle).Methods("GET")
	rt.HandleFunc("/api/articles", CreateArticle).Methods("POST")
	fmt.Println("scraper-api started. Listening on port:8081")
	log.Fatal(http.ListenAndServe(":8081", rt))
}

func startApi(){
	fmt.Println("Starting db connection")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting scraper-api")
	RefreshSources()
	RefreshArticles()
	handleRequests()
}

