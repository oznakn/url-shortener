package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var password string

func decodeHandler(response http.ResponseWriter, request *http.Request, db Database) {
	short := mux.Vars(request)["short"]

	url, err := db.Get(short)

	if err != nil {
		http.Error(response, `{"error": "No such URL"}`, http.StatusNotFound)
		return
	}

	http.Redirect(response, request, url, 302)
}

func encodeHandler(response http.ResponseWriter, request *http.Request, db Database, baseURL string) {
	decoder := json.NewDecoder(request.Body)

	var data struct {
		URL      string `json:"url"`
		Password string `json:"pass"`
	}

	err := decoder.Decode(&data)
	if err != nil || data.Password != password {
		http.Error(response, `{"error": "Unable to parse json"}`, http.StatusBadRequest)
		return
	}

	if !govalidator.IsURL(data.URL) {
		http.Error(response, `{"error": "Not a valid URL"}`, http.StatusBadRequest)
		return
	}

	short, err := db.Save(data.URL)
	if err != nil {
		log.Println(err)
		return
	}

	resp := map[string]string{"url": baseURL + short, "short": short, "error": ""}
	jsonData, _ := json.Marshal(resp)
	response.Write(jsonData)

}

func main() {

	if os.Getenv("PASSWORD") == "" {
		log.Fatal("PASSWORD environment variable must be set")
	}

	password = os.Getenv("PASSWORD")

	if os.Getenv("BASE_URL") == "" {
		log.Fatal("BASE_URL environment variable must be set")
	}

	if os.Getenv("DB_PATH") == "" {
		log.Fatal("DB_PATH environment variable must be set")
	}
	db := sqlite{Path: path.Join(os.Getenv("DB_PATH"), "db.sqlite")}
	db.Init()

	baseURL := os.Getenv("BASE_URL")
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	r := mux.NewRouter()
	r.HandleFunc("/save", func(response http.ResponseWriter, request *http.Request) {
		encodeHandler(response, request, db, baseURL)
	}).Methods("POST")

	r.HandleFunc("/{short}", func(response http.ResponseWriter, request *http.Request) {
		decodeHandler(response, request, db)
	})

	r.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "http://cclub.metu.edu.tr", 301)
	})

	log.Printf("Starting server on port 1337 with password: %s\n", password)
	log.Fatal(http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, r)))
}
