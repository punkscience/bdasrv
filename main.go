package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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

// HomePageParams An object to carry the parameters required to our home page template
type HomePageParams struct {
	Name  string
	URL   string
	Track int
}

// DB Create a global database to hold all the show data
var DB DnBDB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("web/*.htmlgo"))
}

func getHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	min := 0
	max := len(DB.Files)

	track := rand.Intn(max-min) + min

	params := HomePageParams{Track: track, Name: DB.Files[track].Filename, URL: DB.Files[track].URL}

	tpl.ExecuteTemplate(w, "index.htmlgo", params)

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
	fs := http.FileServer(http.Dir("web/"))

	r.HandleFunc("/", getHome).Methods((http.MethodGet))

	// handle static files too
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/show/{showID}", getShow).Methods(http.MethodGet)
	api.HandleFunc("/shows", getShows).Methods(http.MethodGet)

	// If there is a cert file, use it
	if _, err := os.Stat("full-cert.crt"); err == nil {
		log.Fatal(http.ListenAndServeTLS(":8080", "full-cert.crt", "private-key.key", r))
	} else {
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}
