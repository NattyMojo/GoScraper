package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Article struct {
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func AllArticles (w http.ResponseWriter, r *http.Request){
	articles := Articles{
		Article{Title:"Test Title", Desc: "Test Description", Content: "Hello World"},
	}
	fmt.Println("Endpoint Hit: All Articles Endpoint")
	json.NewEncoder(w).Encode(articles)
}

func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}

/*
func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}

func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}

func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}

func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}

func NewArticle(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "New Article endpoint hit.")
}*/