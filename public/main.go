package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Note different representation for struct fields and JSON elements, differentiated by lower-case start letter
type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdOn"`
}

// storage for the Note collection since there are no database.
var noteStore = make(map[string]Note)

// variable to generate key for the collection
var id int = 0

// PostNoteHandler HTTP Post - /api/notes
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	// decode the incoming Note JSON. NewDecoder creates a Decoder object
	// and its Decode method decodes the JSON string into the given type (Note type in this example)
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedOn = time.Now()
	id++
	// the string type is used as the key for the noteStore map, so the int type is converted to string
	// using the strconv.Itoa()
	k := strconv.Itoa(id)
	log.Printf("id: %s", k)
	// add the new note to noteStore map with key = string(id)
	noteStore[k] = note

	// encode the new note JSON to write it in the response
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
	// change the datatype of noteStore to a JSON datatype, e.g. slice.
	for _, v := range noteStore {
		notes = append(notes, v)
	}
	w.Header().Set("Content-Type", "application/json")
	// encode created slice as JSON, by using json.Marshall()
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	// write marshalled data (JSON) to response
	_, err = w.Write(j)
	if err != nil {
		return
	}
}

// PutNoteHandler HTTP Put - /api/notes/{id}
func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	// endpoint is /api/notes/{id}. to get id, use mux.Vars(r). it returns [string]string map
	vars := mux.Vars(r)
	k := vars["id"]
	var noteToUpd Note
	// decode the incoming Note json to noteToUpd
	err = json.NewDecoder(r.Body).Decode(&noteToUpd)
	if err != nil {
		panic(err)
	}
	if note, ok := noteStore[k]; ok {
		// keep createdOn info
		noteToUpd.CreatedOn = note.CreatedOn
		// delete existing item and add the updated item
		delete(noteStore, k)
		noteStore[k] = noteToUpd
	} else {
		log.Printf("could not find key of Note %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteNoteHandler HTTP Delete - /api/notes/{id}
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
