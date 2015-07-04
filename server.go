package main

import (
	"fmt"
	"net/http"
	"html/template"
)


type HomePage struct {
	Title string
	Body  string
}

var page_templates = template.Must(template.ParseFiles(
	"./tpl/head.html",
	"./tpl/page_body.html"))

func loadHomePage(w http.ResponseWriter, r *http.Request) {
	p := &HomePage{Title: "Title", Body: "Body"}
	renderTpl(w, "head", p)
	renderTpl(w, "page_body", p)
}

func renderTpl(w http.ResponseWriter, tmpl string, p *HomePage) {
	err := page_templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.HandleFunc("/", loadHomePage)
	fmt.Printf("tested!!");
	http.ListenAndServe(":8080", nil)
}