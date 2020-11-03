package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const PORT = 8080

type AboutMe struct {
	FullName string `json:"full_name"`
	Fact     string `json:"string"`
}

var list []AboutMe

func getAboutMe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		// serialize a random response
		i := rand.Intn(len(list))
		str, err := json.Marshal(list[i])
		if err != nil {
			http.Error(w, "internal server error on serialization", http.StatusInternalServerError)
			log.Printf("SERIALIZATION ERROR: %v\n", err)
			return
		}

		// write response
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(str)
		if err != nil {
			http.Error(w, "internal server error on write response", http.StatusInternalServerError)
			log.Printf("WRITE ERROR: %v\n", err)
			return
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	// add facts
	name := "Mickey Mouse"
	list = append(list, AboutMe{
		FullName: name,
		Fact:     "I am a mouse.",
	})
	list = append(list, AboutMe{
		FullName: name,
		Fact:     "I am 92 years old.",
	})

	// start listening for requests
	log.Printf("starting server on port %v...\n", PORT)
	mux := http.NewServeMux()
	mux.HandleFunc("/me", getAboutMe)
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), mux)
	if err != nil {
		panic(err)
	}

}
