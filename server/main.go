package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading env file")
	}

}
func main() {

	http.Handle("/", http.FileServer(http.Dir(os.Getenv("STATIC_DIR"))))
	http.HandleFunc("/search", HandleFindPlaces)
	http.HandleFunc("/ping", HandlePing)
	http.HandleFunc("/key", HandleGetKey)
	http.ListenAndServe(":8080", nil)
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"statusCode": 200,
		"data":       "pong",
	})

}

func HandleFindPlaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query()

	location := params.Get("latlng")
	keyword := strings.Replace(params.Get("term"), " ", "+", -1)

	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	fields := "photos,formatted_address,name"

	base := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?fields"
	url := fmt.Sprintf("%s=%s&location=%s&radius=150000&keyword=%s&key=%s", base, fields, location, keyword, key)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func HandleGetKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"statusCode": 200,
		"key":        key,
	})

}
