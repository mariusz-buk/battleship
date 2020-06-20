package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	mux "github.com/gorilla/mux"
	sessions "github.com/gorilla/sessions"
)

type command struct {
	name string
	fn   func(http.ResponseWriter, *http.Request)
}

type message struct {
	Command string `json:"command"`
	Message string `json:"message"`
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(".8pqDQ{%&aWE:hDP]Jv5hF<")
	store = sessions.NewCookieStore(key)
)

func addTemplateIterator() template.FuncMap {
	return template.FuncMap{
		"Iterate": func(count uint) []uint {
			var i uint
			var Items []uint
			for i = 1; i <= count; i++ {
				Items = append(Items, i)
			}
			return Items
		},
	}
}

func checkSession(writer http.ResponseWriter, request *http.Request) *sessions.Session {
	session, _ := store.Get(request, "session")

	if session.IsNew {
		session.Values["step"] = 1
		session.Save(request, writer)
	}

	return session
}

func getIndexFile(writer http.ResponseWriter, request *http.Request) {
	checkSession(writer, request)

	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}
	tpl := template.New("battleship")
	tpl.Funcs(addTemplateIterator())
	tpl, _ = tpl.ParseFiles(allFiles...)
	err = tpl.ExecuteTemplate(writer, "index.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func step1(writer http.ResponseWriter, request *http.Request) {
	msg := message{
		Command: "display",
		Message: "Randomly setting ships position",
	}
	json.NewEncoder(writer).Encode(msg)
}

func getStep(writer http.ResponseWriter, request *http.Request) {
	session := checkSession(writer, request)
	//vars := mux.Vars(request)
	//vars["step"
	var step uint = 1
	if session.Values["step"] != nil {
		step = session.Values["step"].(uint)
	}

	switch step {
	case 1:
		step1(writer, request)
	}

	step += 1
	if step > 1 {
		step = 1
	}

	session.Values["step"] = step
	session.Save(request, writer)
}

func addCommands(r *mux.Router) {
	r.HandleFunc("/", getIndexFile)
	r.HandleFunc("/get-step", getStep)
}

func main() {
	r := mux.NewRouter()
	addCommands(r)

	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))

	log.Fatal(http.ListenAndServe(":80", r))
}
