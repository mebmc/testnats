package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Entry struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Guess     string `json:"guess"`
}

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type apiserver struct {
}

func (apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.RequestURI == "/api/submit" {
		body, _ := ioutil.ReadAll(r.Body)
		data := Entry{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("Err: %v", err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{true, ""})
	}
}

func main() {
	PORT := 8080
	fmt.Printf("Listening on %v\n", PORT)

	api := apiserver{}

	mux := http.NewServeMux()
	mux.Handle("/api/", api)

	LSTR := fmt.Sprintf(":%v", PORT)
	log.Fatal(http.ListenAndServe(LSTR, mux))
}
