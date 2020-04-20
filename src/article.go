package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Article struct {
	ID int `json:"id"`
	Address string `json:"address"`
	Headline string `json:"headline"`
	Source *Source `json:"source"`
}

type Source struct {
	ID int `json:"id"`
	Address string `json:"address"`
	Name string `json:"name"`
}

//This is used as a temp article before we know which source it should map to.
type Tarticle struct {
	Address string `json:"address"`
	Headline string `json:"headline"`
	SourceName string `json:"source"`
}

var (
	sources []Source
	articles []Article
)

func GetSources(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	RefreshSources()
	json.NewEncoder(w).Encode(sources)
}

func RefreshSources(){
	var s []Source
	rows, err := Db.Query("Select * from articles.source")
	defer rows.Close()
	for rows.Next(){
		var source_id int
		var name string
		var address string
		err = rows.Scan(&source_id, &name, &address)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, Source{source_id, address, name})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	sources = nil
	sources = append([]Source(nil), s...) //Populates sources
}


func GetSource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //get params
	w.Header().Set("Content-Type", "application/json")
	s := getSource(params["name"])
	if s.ID != -1 {
		json.NewEncoder(w).Encode(s)
	}
}

func getSource(s string) Source {
	var source Source
	for _, source = range sources {
		if source.Name == s {
			return source
		}
	}
	source.ID = -1 //Set id as bad value
	return source
}


func CreateSource(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var sid int
	var source Source
	_ = json.NewDecoder(r.Body).Decode(&source)
	q := `insert into articles.source(name, address) values ($1, $2) returning source_id;`
	err := Db.QueryRow(q,source.Name,source.Address).Scan(&sid) //returns generated source_Id
	if err != nil {
		fmt.Println("Something went wrong creating the source.")
		return
	}
	source.ID = sid //sticks the source_id on the source struct
	sources = append(sources, source)
	json.NewEncoder(w).Encode(source)
}

func RefreshArticles(){
	var a []Article
	q := `	Select
			article.article_id, 
			article.address, 
			article.headline, 
			source.name 
			from articles.article
			join articles.source on source.source_id = article.source_id`
	rows, err := Db.Query(q)
	defer rows.Close()
	for rows.Next(){
		var artID int
		var address string
		var headline string
		var sName string
		err = rows.Scan(&artID, &address, &headline, &sName)
		if err != nil {
			log.Fatal(err)
		}
		var src Source = getSource(sName)
		var art Article = Article{
			ID:       artID,
			Address:  address,
			Headline: headline,
			Source:   &src,
		}
		a = append(a, art)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	articles = nil
	articles = append([]Article(nil), a...) //Populates sources
}

func GetArticles(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	RefreshArticles()
	json.NewEncoder(w).Encode(articles)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //get params
	w.Header().Set("Content-Type", "application/json")
	a := getArticle(params["address"])
	if a.ID != -1 {
		json.NewEncoder(w).Encode(a)
	}
}

func getArticle(s string) Article{
	var article Article
	for _, at := range articles {
		if at.Address == s {
			return at
		}
	}
	return article
}

func CreateArticle(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var aid int
	var art Article
	var tart Tarticle
	_ = json.NewDecoder(r.Body).Decode(&tart) //decodes the temp article object
	ab := getArticle(tart.Address)
	if ab.ID != 0{
		return
	}
	art.Address = tart.Address
	art.Headline = tart.Headline
	source  := getSource(tart.SourceName) //then grabs the source based on the name from tart.
	art.Source = &source
	q := 	`INSERT INTO articles.article(address, headline, source_id)
			SELECT $1, $2, source.source_id
			FROM articles.source
			WHERE source.name = $3
			returning article_id;`
	err := Db.QueryRow(q, art.Address, art.Headline, source.Name).Scan(&aid) //returns generated article
	if err != nil {
		fmt.Println("Something went wrong creating the article.")
		fmt.Println(err)
		return
	}
	art.ID = aid //sticks the articleid on the art struct
	articles = append(articles, art)
	json.NewEncoder(w).Encode(art)
}