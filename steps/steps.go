package steps

import (
	"encoding/json"
	"net/http"

	"battleship.go/sessions"
	"battleship.go/ships"
)

func firstMessage(writer http.ResponseWriter, request *http.Request) {
	msg := struct {
		Command string
		Message string
	}{
		Command: "display",
		Message: "Randomly setting ships position",
	}
	json.NewEncoder(writer).Encode(msg)
}

func prepareShips(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)

	if ships.Init(writer, request) {
		msg := struct {
			Command string
			Boards  ships.BlackAndWhiteArmies
		}{
			Command: "fill-boards",
			Boards:  session.Values["ships"].(ships.BlackAndWhiteArmies),
		}
		json.NewEncoder(writer).Encode(msg)
	} else {
		msg := struct {
			Command string
			Message string
		}{
			Command: "display",
			Message: "Could not place all ships on the boards in 100 rounds",
		}
		json.NewEncoder(writer).Encode(msg)
	}
}

func theEnd(writer http.ResponseWriter, request *http.Request) {
	msg := struct {
		Command string
		Message string
	}{
		Command: "the-end",
		Message: "Demo is over",
	}
	json.NewEncoder(writer).Encode(msg)
}

func ProceedWithSteps(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)
	step := 1
	if session.Values["step"] != nil {
		step = session.Values["step"].(int)
	}

	switch step {
	case 1:
		firstMessage(writer, request)
	case 2:
		prepareShips(writer, request)
	default:
		theEnd(writer, request)
	}

	step++
	session.Values["step"] = step
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
