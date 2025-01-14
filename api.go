package main

import (
	"encoding/json"
	"strings"
	//"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		Error string `json:"error"`
	}
	thisError := errorStruct{Error: msg}
	res, _ := json.Marshal(thisError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func healthResponseHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func validationResponseHandler(w http.ResponseWriter, req *http.Request) {
	type validationRequest struct {
		Body string `json:"body"`
	}
	type validationResponse struct {
		Cleaned_body string `json:"cleaned_body,omitempty"`
	}
	thisRequest := validationRequest{}
	thisResponse := validationResponse{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)

	if err != nil {
		respondWithError(w, 400, "Something went wrong")
	}

	if len(thisRequest.Body) >= 140 {
		respondWithError(w, 400, "Chirpy too long")
	}

	cleanedBody, _ := cleanBody(thisRequest.Body)
	thisResponse.Cleaned_body = cleanedBody
	respondWithJSON(w, 200, thisResponse)
}

func cleanBody(body string) (cleaned_body string, was_cleaned bool) {
	wordBlacklist := []string{"Kerfuffle", "Sharbert", "Fornax"}

	split := strings.Split(body, " ")
	for i, word := range split {
		for _, badWord := range wordBlacklist {
			if word == badWord {
				split[i] = "****"
			}
		}
		for _, badWord := range wordBlacklist {
			if word == strings.ToLower(badWord) {
				split[i] = "****"
			}
		}
	}
	cleaned_body = strings.Join(split, " ")
	was_cleaned = !(cleaned_body == body)

	return cleaned_body, was_cleaned
}
