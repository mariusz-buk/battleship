package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type command struct {
	name string
	fn   func()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func getConnectionToken() {

}

func main() {
	commands := []command{
		{
			name: "getConnectionToken",
			fn:   getConnectionToken,
		},
	}
	for _, command := range commands {
		http.HandleFunc("/"+command.name, command.fn())
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}
