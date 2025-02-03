package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)


func (cfg *apiConfig)  healthHandler(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.env != "dev"{
		w.WriteHeader(http.StatusForbidden)
	}
	cfg.fileserverHits.Store(0)
	err:=cfg.db.DeleteAllUsers(r.Context())
	if err!=nil{
		ServerErrorResponse(w)
	}
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




// USER HANDLERS


func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	var output struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	var errr errror
	w.Header().Set("Content-Type","application/json")
	err:= json.NewDecoder(r.Body).Decode(&input)
	if err!=nil{
		badRequestErrorResponse(w,http.StatusBadRequest,errr)
	}
	slog.Info("email",input.Email)
	user,err:= cfg.db.CreateUser(r.Context(),input.Email)
	if err!=nil{
		ServerErrorResponse(w)
	}
	output.ID = user.ID
	output.CreatedAt = user.CreatedAt.Time
	output.UpdatedAt = user.UpdatedAt.Time
	output.Email = user.Email
	slog.Info("user",user)

	data, err:= json.Marshal(output)
	w.WriteHeader(http.StatusCreated)
	w.Write(data)


}






//Chirp HANDLERS


func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Body string `json:"body"`
	}
	var er struct{
		Error  string `json:"error"`
	}
	var success struct {
		Valid bool
	}	
	var cleanedBody struct {
		Cleaned_Body string `json:"cleaned_body"`
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