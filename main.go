package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Show Structure to hold the data about a particular show
type Show struct {
	Event      string
	URL        string
	Filename   string
	Downloaded bool
}

// DnBDB define data structure
type DnBDB struct {
	Output   string
	LastScan string
	Files    []Show
}

// ShowResponse A response object containing an array of a list of all available shows
type ShowResponse struct {
	Shows []string `json:"shows"`
}

// DB Create a global database to hold all the show data
var DB DnBDB

func getShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Get a complete list of shows
	showResp := new(ShowResponse)

	// Loop through all the events and get the show names
	for _, show := range DB.Files {
		showResp.Shows = append(showResp.Shows, show.Event)
	}

	// Convert the object to JSON
	bytes, err := json.MarshalIndent(showResp, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Print(string(bytes))
	w.Write(bytes)

}

func getHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
}

func main() {
	// Read in the database
	data, err := ioutil.ReadFile("./filedata.json")
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data, &DB)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Loaded %d shows and ready for work.\n", len(DB.Files))

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/", getHome).Methods(http.MethodGet)

	api.HandleFunc("/shows", getShows).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)

	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
