package main

import (
	"html/template"
	"log"
	"net/http"

	"os"

	"github.com/lutepluto/butler/wechat"
)

var templates = template.Must(template.ParseFiles("index.html"))

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[DEBUG]")
	log.SetOutput(os.Stdout)
}

func main() {
	log.Printf("Start serving...")
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, request *http.Request) {
	session := wechat.DefaultSession

	if err := templates.ExecuteTemplate(w, "index.html", session.UUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
