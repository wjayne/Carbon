package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"fmt"
	"regexp"
	"errors"
)
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct{
	Title string
	Body []byte
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error){
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil{
		http.NotFound(w,r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2],nil
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body,err :=ioutil.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	//fmt.Println(title, "HERE")
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title,err := getTitle(w,r)
	if err != nil{
		return
	}
	p,err :=loadPage(title)
	if err != nil{
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	fmt.Println(p, "Page")
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title, err := getTitle(w,r)
	if err != nil{
		return
	}
	p,err := loadPage(title)
	if err != nil{
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title,err := getTitle(w,r)
	if err != nil{
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", viewHandler)
	http.ListenAndServe(":8080",nil)
}