package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func main() {
	log.Println("Starting Windowpane")
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("files/"))

	mux.Handle("/files/", http.StripPrefix("/files", fileServer))
	mux.HandleFunc("/", getTemplate)

	checkError(http.ListenAndServe(":8080", mux))
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getTemplate(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.html")

	c := getClient()
	user, resp := getBasicInfo(c)
	if resp.StatusCode == http.StatusForbidden {
		http.Error(w, "403 Forbidden (likely rate limited)", http.StatusForbidden)
		return
	} else if resp.StatusCode != 200 {
		log.Println(resp.Body)
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	repos, resp := getRepositories(c)
	if resp.StatusCode == http.StatusForbidden {
		http.Error(w, "403 Forbidden (likely rate limited)", http.StatusForbidden)
		return
	} else if resp.StatusCode != 200 {
		log.Println(resp.Body)
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	events, resp := getEvents(c)
	if resp.StatusCode == http.StatusForbidden {
		http.Error(w, "403 Forbidden (likely rate limited)", http.StatusForbidden)
		return
	} else if resp.StatusCode != 200 {
		log.Println(resp.Body)
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	data := ghData{
		UserInfo: user,
		Repos:    repos,
		Events:   events,
	}

	tmpl, err := template.ParseGlob(layout)
	checkError(err)

	err = tmpl.ExecuteTemplate(w, "layout", data)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
