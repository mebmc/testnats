package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	stan "github.com/nats-io/go-nats-streaming"
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
	Sc stan.Conn
}

func (a apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method == "POST" && r.RequestURI == "/api/submit" {
		body, _ := ioutil.ReadAll(r.Body)
		data := Entry{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("Err: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{false, "invalid data submited"})
		}

		err = a.Sc.Publish(data.ID, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{false, "unable to publish"})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{true, ""})
	}
}

func main() {
	PORT := 8080
	URL := "nats://docker.local:4222"
	fmt.Printf("Listening on %v\n", PORT)
	// fmt.Println(stan.DefaultNatsURL)

	sc, err := stan.Connect("test-cluster", "api", stan.NatsURL(URL))
	if err != nil {
		log.Printf("Err: %v\n", err)
	}

	api := apiserver{Sc: sc}

	mux := http.NewServeMux()
	mux.Handle("/api/", api)

	LSTR := fmt.Sprintf(":%v", PORT)
	log.Fatal(http.ListenAndServe(LSTR, mux))
}
