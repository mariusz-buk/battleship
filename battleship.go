package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"battleship.go/sessions"
	"battleship.go/steps"
	"github.com/gorilla/mux"
)

type boardData struct {
	TopArray  [10]string
	LeftArray [11]string
}

func getBoardData() boardData {
	var board boardData
	board.LeftArray[0] = ""
	for i := 0; i < 10; i++ {
		board.TopArray[i] = strconv.Itoa(i + 1)
		board.LeftArray[i+1] = string(int('A') + i)
	}
	return board
}

func getIndexFile(writer http.ResponseWriter, request *http.Request) {
	sessions.ClearSession(writer, request)

	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}
	tpl := template.New("battleship")
	tpl, err = template.ParseFiles(allFiles...)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tpl.ExecuteTemplate(writer, "index.html", getBoardData())
	if err != nil {
		fmt.Println(err)
	}
}

func addCommands(r *mux.Router) {
	r.HandleFunc("/", getIndexFile)
	r.HandleFunc("/get-step", steps.ProceedWithSteps)
}

func main() {
	r := mux.NewRouter()
	addCommands(r)

	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))

	var err error
	for port := 80; port < 90; port++ {
		fmt.Printf("\nTrying port %d ...", port)
		err = http.ListenAndServe(":"+strconv.Itoa(port), r)
	}

	log.Fatal(err)
}
