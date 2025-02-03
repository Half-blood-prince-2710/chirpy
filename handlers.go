package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

)


func (cfg *apiConfig)  healthHandler(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")

	html:= `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	htmlcontent:= fmt.Sprintf(html,cfg.fileserverHits.Load())
	fmt.Fprint(w,htmlcontent)
}


func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Body string `json="body"`
	}
	var er struct{
		Error  string `json="error"`
	}
	var success struct {
		Valid bool
	}	
	var cleanedBody struct {
		Cleaned_Body string `json="cleaned_body"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err!=nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-type","application/json")
		er.Error = "Something went wrong"
		dat,err := json.Marshal(er)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		w.Write(dat)
		return
	}
	// fmt.Print(len(input.Body),"\n","body",input.Body,"\n")
	if len(input.Body) >140 {
			w.WriteHeader(400)
		w.Header().Set("Content-type","application/json")
		er.Error = "Chirp is too long"
		dat,err := json.Marshal(er)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		w.Write(dat)
		return
	}
	bannedWords:= []string{"kerfuffle","sharbert","fornax"}
	for _,word := range bannedWords {
		re := regexp.MustCompile(`\b(?i)` + word + `\b`)
		input.Body=re.ReplaceAllString(input.Body,"****")
	}
	cleanedBody.Cleaned_Body = input.Body
	w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type","application/json")
		success.Valid = true
		dat,err := json.Marshal(cleanedBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		
	w.Write(dat)

}