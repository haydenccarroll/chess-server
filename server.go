// The following guide helped me greatly in creating this basic server.
// https://andela.com/insights/using-golang-to-create-a-restful-json-api/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	chess "github.com/malbrecht/chess"
)

type move struct {
	Time time.Time `json:"Time"`
	Move string `json:"Move"`
}

type chessBoardState struct {
	Moves       []move `json:"Moves"`
	Fen         string `json:"Fen"`
	IsWhiteMove bool   `json:"IsWhiteMove"`
	MoveCount   int    `json:"MoveCount"`
}

var boardState = chessBoardState{
	Moves:       nil,
	Fen:         "",
	IsWhiteMove: true,
	MoveCount:   0,
}

var myBoard = new(chess.Board)

func createBoard() {
	for i := 0; i < 8; i++ {
		myBoard.Piece[i+8] = chess.WP
		myBoard.Piece[63-8-i] = chess.BP
	}

	myBoard.Piece[0] = chess.WR
	myBoard.Piece[1] = chess.WN
	myBoard.Piece[2] = chess.WB
	myBoard.Piece[3] = chess.WQ
	myBoard.Piece[4] = chess.WK
	myBoard.Piece[5] = chess.WB
	myBoard.Piece[6] = chess.WN
	myBoard.Piece[7] = chess.WR

	myBoard.Piece[63-0] = chess.BR
	myBoard.Piece[63-1] = chess.BN
	myBoard.Piece[63-2] = chess.BB
	myBoard.Piece[63-3] = chess.BQ
	myBoard.Piece[63-4] = chess.BK
	myBoard.Piece[63-5] = chess.BB
	myBoard.Piece[63-6] = chess.BN
	myBoard.Piece[63-7] = chess.BR
}

func createMove(w http.ResponseWriter, r *http.Request) {
	var newMove move
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please ensure the move has a timestamp and a move")
	}

	json.Unmarshal(reqBody, &newMove)
	var currentMove, error = myBoard.ParseMove((newMove.Move))
	if error != nil {
		fmt.Fprintf(w, "The move could not be played. Are you sure you formatted it correctly?")
	} else {
		t := time.Now()
		var structMoveNew move
		structMoveNew.Move = newMove.Move
		structMoveNew.Time = t
		myBoard = myBoard.MakeMove(currentMove)
		boardState.Moves = append(boardState.Moves, structMoveNew)
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(newMove)
}

// func getOneEvent(w http.ResponseWriter, r *http.Request) {
// 	eventID := mux.Vars(r)["id"]

// 	for _, singleEvent := range events {
// 		if singleEvent.ID == eventID {
// 			json.NewEncoder(w).Encode(singleEvent)
// 		}
// 	}
// }

func getBoardState(w http.ResponseWriter, r *http.Request) {
	boardState.Fen = myBoard.Fen()
	if myBoard.SideToMove == 0 {
		boardState.IsWhiteMove = true
	} else {
		boardState.IsWhiteMove = false
	}
	boardState.MoveCount = myBoard.MoveNr
	json.NewEncoder(w).Encode(boardState)
}

// func updateEvent(w http.ResponseWriter, r *http.Request) {
// 	eventID := mux.Vars(r)["id"]
// 	var updatedEvent event

// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
// 	}
// 	json.Unmarshal(reqBody, &updatedEvent)

// 	for i, singleEvent := range events {
// 		if singleEvent.ID == eventID {
// 			singleEvent.Title = updatedEvent.Title
// 			singleEvent.Description = updatedEvent.Description
// 			events = append(events[:i], singleEvent)
// 			json.NewEncoder(w).Encode(singleEvent)
// 		}
// 	}
// }

// func deleteEvent(w http.ResponseWriter, r *http.Request) {
// 	eventID := mux.Vars(r)["id"]

// 	for i, singleEvent := range events {
// 		if singleEvent.ID == eventID {
// 			events = append(events[:i], events[i+1:]...)
// 			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
// 		}
// 	}
// }

func main() {
	createBoard()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", createMove).Methods("POST")
	router.HandleFunc("/", getBoardState).Methods("GET")
	// router.HandleFunc("/", homeLink)
	// router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	// router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	// router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
