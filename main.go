package main

import (
	"encoding/json"
	"io/ioutil"
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

func getVersion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		// write response
		w.Header().Add("Content-Type", "test/plain")
		str := []byte("1.0.0")
		_, err := w.Write(str)
		if err != nil {
			http.Error(w, "internal server error on write response", http.StatusInternalServerError)
			log.Printf("WRITE ERROR: %v\n", err)
			return
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

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

func getAboutOther(w http.ResponseWriter, r *http.Request, other string, header string) {
	switch r.Method {
	case http.MethodGet:

		// create client
		client := &http.Client{}

		// create request
		req, err := http.NewRequest("GET", "http://"+other+"/me", nil)
		if err != nil {
			http.Error(w, "internal server error on create request", http.StatusInternalServerError)
			log.Printf("CREATE ERROR: %v\n", err)
			return
		}

		// add the routing header
		if len(header) > 0 {
			req.Header.Add("kubernetes-route-as", header)
		}

		// fetch the request
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "internal server error on GET mickey", http.StatusInternalServerError)
			log.Printf("GET ERROR: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "internal server error on read response", http.StatusInternalServerError)
			log.Printf("RESPONSE ERROR: %v\n", err)
			return
		}

		// write response
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			http.Error(w, "internal server error on write response", http.StatusInternalServerError)
			log.Printf("WRITE ERROR: %v\n", err)
			return
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAboutMickey(w http.ResponseWriter, r *http.Request) {
	getAboutOther(w, r, "mickey", "")
}

func getAboutLila(w http.ResponseWriter, r *http.Request) {
	getAboutOther(w, r, "lila", "evmolyva-e52d")
}

func main() {

	// add facts
	name := "Peter Lasne"
	list = append(list, AboutMe{
		FullName: name,
		Fact:     "I have a 3-year old.",
	})
	list = append(list, AboutMe{
		FullName: name,
		Fact:     "I am an avid table-top gamer.",
	})

	// start listening for requests
	log.Printf("starting server on port %v...\n", PORT)
	mux := http.NewServeMux()
	mux.HandleFunc("/me", getAboutMe)
	mux.HandleFunc("/mickey", getAboutMickey)
	mux.HandleFunc("/lila", getAboutLila)
	mux.HandleFunc("/version", getVersion)
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), mux)
	if err != nil {
		panic(err)
	}

}
