package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type command struct {
	name string
	fn   func(http.ResponseWriter, *http.Request)
}

type player struct {
	token string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func addTemplateIterator() {
	template.FuncMap{
		"Iterate": func(count *uint) []uint {
			var i uint
			var Items []uint
			for i = 0; i < (*count); i++ {
				Items = append(Items, i)
			}
			return Items
		},
	}
}

func getConnectionToken(writer http.ResponseWriter, request *http.Request) {

}

func addCommands() {
	commands := []command{
		{
			name: "getConnectionToken",
			fn:   getConnectionToken,
		},
	}
	for _, command := range commands {
		http.HandleFunc("/"+command.name, command.fn)
	}
}

func main() {
	addTemplateIterator()
	addCommands()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
