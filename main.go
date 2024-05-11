package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func main() {
	shorten := &ShortenerUrl{
		url: make(map[string]string),
	}

	http.HandleFunc("/shorten", shorten.ShorterLinkHandler)
	http.HandleFunc("/", shorten.RedirectLinks)

	fmt.Println("Url Api run in PORT 8080")
	http.ListenAndServe(":8080", nil)
}

type ShortenerUrl struct {
	url map[string]string
}

type Link struct {
	Links string `json:"link"`
}

//Util

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

//Handler

func(s *ShortenerUrl) ShorterLinkHandler(w http.ResponseWriter, r *http.Request) {
	sk := generateShortKey()

	var req Link

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	s.url["/"+sk] = req.Links

	shortenLink := fmt.Sprintf("localhost:8080/%s", sk)

	res := Link{
		Links: shortenLink,
	}

	jsonData, err := json.Marshal(res)

	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)
}

func(s *ShortenerUrl) RedirectLinks(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path
	fmt.Println(shortKey, "GINISIH")
	if shortKey == "" {
        http.Error(w, "Missing shortkey", http.StatusBadRequest)
        return
    }

	originalLink, isFound := s.url[shortKey]

	if !isFound {
		http.Error(w, "Shortkey not found", http.StatusNotFound)
        return
	}

	http.Redirect(w, r, originalLink, http.StatusPermanentRedirect)
}



