package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdOn"`
}

// store for the Notes collection
var noteStore = make(map[string]Note)

// variable to generate key for the collection
var id = 0

// PostNoteHandler HTTP Post - /api/notes
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	// decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	log.Printf("id: %s", k)
	noteStore[k] = note

	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(j)
	if err != nil {
		return
	}
}

// GetNoteHandler HTTP Get - /api/notes
func GetNoteHandler(w http.ResponseWriter, _ *http.Request) {
	var notes []Note
	for _, v := range noteStore {
		notes = append(notes, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		return
	}
}

// PutNoteHandler HTTP Put - /api/notes
func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	var noteToUpd Note
	// decode the incoming Note json
	err = json.NewDecoder(r.Body).Decode(&noteToUpd)
	if err != nil {
		panic(err)
	}
	if note, ok := noteStore[k]; ok {
		noteToUpd.CreatedOn = note.CreatedOn
		// delete existing item and add the updated item
		delete(noteStore, k)
		noteStore[k] = noteToUpd
	} else {
		log.Printf("could not find key of Note %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteNoteHandler HTTP Delete - /api/notes
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	// remove from store
	if _, ok := noteStore[k]; ok {
		// delete existing item
		delete(noteStore, k)
	} else {
		log.Printf("could not find key of Note %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

// entry point of the program
func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("listening...")
	err := server.ListenAndServe()
	if err != nil {
		return
	}

}
