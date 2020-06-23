package steps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"battleship.go/sessions"
	"battleship.go/ships"
)

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

func translatePositionX(x int) string {
	return strconv.Itoa(x + 1)
}

func translatePositionY(y int) string {
	return string(int('A') + y)
}

func blackFirePrep(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)

	x, y := ships.GetTarget()

	session.Values["nextTargetX"] = x
	session.Values["nextTargetY"] = y

	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := struct {
		Command, Message, Board string
		X, Y                    int
	}{
		Command: "new-target",
		Message: fmt.Sprintf("Black fires %s:%s", translatePositionY(y), translatePositionX(x)),
		Board:   "white",
		X:       x,
		Y:       y,
	}
	json.NewEncoder(writer).Encode(msg)
}

func whiteFirePrep(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)

	x, y := ships.GetTarget()

	session.Values["nextTargetX"] = x
	session.Values["nextTargetY"] = y

	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := struct {
		Command, Message, Board string
		X, Y                    int
	}{
		Command: "new-target",
		Message: fmt.Sprintf("White fires %s:%s", translatePositionY(y), translatePositionX(x)),
		Board:   "black",
		X:       x,
		Y:       y,
	}
	json.NewEncoder(writer).Encode(msg)
}

func blackFires(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)

	armies := session.Values["ships"].(ships.BlackAndWhiteArmies)
	x := session.Values["nextTargetX"].(int)
	y := session.Values["nextTargetY"].(int)
	report := armies.Fire("black", x, y)

	msg := struct {
		Command, Message, TranslatedX, TranslatedY, Board string
		X, Y                                              int
		Hit, Sank                                         bool
	}{
		Command:     "hit-report",
		Message:     fmt.Sprintf("Black fires %s:%s", translatePositionY(report.Y), translatePositionX(report.X)),
		X:           report.X,
		Y:           report.Y,
		TranslatedX: translatePositionX(report.X),
		TranslatedY: translatePositionY(report.Y),
		Hit:         report.Hit,
		Sank:        report.Sank,
		Board:       "white",
	}
	json.NewEncoder(writer).Encode(msg)
}

func whiteFires(writer http.ResponseWriter, request *http.Request) {
	session := sessions.CheckSession(writer, request)

	armies := session.Values["ships"].(ships.BlackAndWhiteArmies)
	x := session.Values["nextTargetX"].(int)
	y := session.Values["nextTargetY"].(int)
	report := armies.Fire("white", x, y)

	msg := struct {
		Command, Message, TranslatedX, TranslatedY, Board string
		X, Y                                              int
		Hit, Sank                                         bool
	}{
		Command:     "hit-report",
		X:           report.X,
		Y:           report.Y,
		TranslatedX: translatePositionX(report.X),
		TranslatedY: translatePositionY(report.Y),
		Hit:         report.Hit,
		Sank:        report.Sank,
		Board:       "black",
	}
	json.NewEncoder(writer).Encode(msg)
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
		prepareShips(writer, request)
	case 2:
		blackFirePrep(writer, request)
	case 3:
		blackFires(writer, request)
	case 4:
		whiteFirePrep(writer, request)
	case 5:
		whiteFires(writer, request)
		step = 1
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
