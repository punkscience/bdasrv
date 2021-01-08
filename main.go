package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// ShowsResponse A response object containing an array of a list of all available shows
type ShowsResponse struct {
	Shows []string `json:"shows"`
}

// ShowResponse A response object containing the URL for a requested show
type ShowResponse struct {
	URL string `json:"url"`
}

// DB Create a global database to hold all the show data
var DB DnBDB

func getHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	min := 0
	max := len(DB.Files)

	track := rand.Intn(max-min) + min

	msg := "<html><h1>Bassdrive Archive Randomizer</h1><p>This is presently just a server sandbox for experiments. There is nothing to see here.</p>"
	file := DB.Files[track].Filename
	url := DB.Files[track].URL

	fmt.Fprintf(w, msg)
	fmt.Fprintf(w, "<p>Brought to you by <a href=\"http://twitter.com/realpunkscience\">Darryl Wright</a> using the amazing and generous archive of <a href=\"http://bassdrive.com\">Bassdrive.com</a>.</p>")
	fmt.Fprintf(w, "<p>%s</p>", file)
	fmt.Fprintf(w, "<p><audio controls><source src=\"%s\" type=\"audio/mpeg\"></audio></p>", url)
	fmt.Fprintf(w, "</html>")
}

func getShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Get a complete list of shows
	showResp := new(ShowsResponse)

	// Loop through all the events and get the show names
	for _, show := range DB.Files {
		showResp.Shows = append(showResp.Shows, show.Filename)
	}

	// Convert the object to JSON
	bytes, err := json.MarshalIndent(showResp, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Print(string(bytes))
	w.Write(bytes)

}

func getShow(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if val, ok := pathParams["showID"]; ok {
		showID, err := strconv.Atoi(val)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Looking up record %d...\n", showID)

		// Loop through all the events and get the show names
		resp := new(ShowResponse)
		resp.URL = DB.Files[showID].URL

		bytes, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(bytes)
	}
}

func main() {
	// Seed the randomizer
	rand.Seed(time.Now().UnixNano())

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

	r.HandleFunc("/", getHome).Methods((http.MethodGet))

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/show/{showID}", getShow).Methods(http.MethodGet)
	api.HandleFunc("/shows", getShows).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
