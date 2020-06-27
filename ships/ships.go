package ships

import (
	"encoding/gob"
	"math/rand"
	"net/http"

	"battleship.go/sessions"
)

const fromLeftToRight = 0
const fromTopToBottom = 1
const boardSize = 10
const freeSpaceIndex = 99

var shipSizes = [8]int{5, 4, 3, 3, 2, 2, 1, 1}

func Init(writer http.ResponseWriter, request *http.Request) bool {
	gob.Register(BlackAndWhiteArmies{})

	session := sessions.CheckSession(writer, request)
	ships, success := initShips()
	if !success {
		return false
	}

	session.Values["ships"] = ships
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func initShips() (BlackAndWhiteArmies, bool) {
	white := shipsArmy{}
	black := shipsArmy{}
	if !white.setRandomPositions() {
		return BlackAndWhiteArmies{}, false
	}
	if !black.setRandomPositions() {
		return BlackAndWhiteArmies{}, false
	}
	return BlackAndWhiteArmies{
		Black: black,
		White: white,
	}, true
}

func GetTarget() (x, y int) {
	return rand.Intn(boardSize), rand.Intn(boardSize)
}
