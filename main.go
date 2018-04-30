package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// Type that represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	// Create a new chat room
	r := newRoom()
	// Root path
	chatTemplate := &templateHandler{filename: "chat.html"}
	http.HandleFunc("/", chatTemplate.ServeHTTP)
	http.HandleFunc("/room", r.ServeHTTP)
	// get the room going
	go r.run()
	// start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
