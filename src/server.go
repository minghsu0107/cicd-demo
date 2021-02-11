package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Pair is the http request type
type Pair struct {
	A int `json:"a"`
	B int `json:"b"`
}

// AddResult is the http response type
type AddResult struct {
	Sum int `json:"sum"`
}

type addHandler struct{}

func (h *addHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var pair Pair
	err := json.NewDecoder(req.Body).Decode(&pair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sum := pair.A + pair.B
	result := &AddResult{sum}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	port := os.Getenv("PORT")
	http.Handle("/sum", &addHandler{})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
