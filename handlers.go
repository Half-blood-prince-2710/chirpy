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
	"github.com/half-blood-prince-2710/chirpy/internal/auth"
	"github.com/half-blood-prince-2710/chirpy/internal/database"
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
		Password string `json:"password"`
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
	hash,err:=auth.HashPassword(input.Password)
	if err!=nil{
		ServerErrorResponse(w)
	}
	dat:=database.CreateUserParams{
		Email: input.Email,
		HashedPassword: hash,
	}
	slog.Info("email",input.Email)
	user,err:= cfg.db.CreateUser(r.Context(),dat)
	if err!=nil{
		ServerErrorResponse(w)
	}
	output.ID = user.ID
	output.CreatedAt = user.CreatedAt
	output.UpdatedAt = user.UpdatedAt
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
		UserId uuid.UUID `json:"user_id"`
	}
	var er struct{
		Error  string `json:"error"`
	}
	var success struct {
		Valid bool
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
	ch:=database.CreateChirpParams{
		Body: input.Body,
		UserID: input.UserId,
	}
	slog.Info("chirp",ch)
	chirp , err:= cfg.db.CreateChirp(r.Context(),ch)
	if err!=nil{
		slog.Error("err: ",err)
		ServerErrorResponse(w)
		return
	}
	var output struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	output.Body = chirp.Body
	output.ID = chirp.ID
	output.CreatedAt = chirp.CreatedAt
	output.UpdatedAt = chirp.UpdatedAt
	output.UserID = chirp.UserID
	w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-type","application/json")
		success.Valid = true
		dat,err := json.Marshal(output)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		
	w.Write(dat)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps,err:= cfg.db.GetAllChirps(r.Context())
	dbErrorReponse(err,w)

	w.WriteHeader(http.StatusOK)
	data,err := json.Marshal(chirps)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
	w.Write(data)
}


func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id:=r.PathValue("id")
	idx, err:=uuid.Parse(id)
	if err!=nil{
		slog.Error("err","err parsing uuid")
		ServerErrorResponse(w)
	}
	chirp,err:= cfg.db.GetOneChirp(r.Context(),idx)
	slog.Info("chirp",chirp)
	dbErrorReponse(err,w)

	w.WriteHeader(http.StatusOK)
	data,err := json.Marshal(chirp)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
	w.Write(data)
}




































// AUTH HANDLERS

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type output struct {
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
	user,err:=cfg.db.FindUserByEmail(r.Context(),input.Email)
	if err!=nil {
		unauthorizedErrorResponse(errr,w)
	}

	err =auth.CheckPasswordHash(input.Password,user.HashedPassword)
	if err!=nil {
		unauthorizedErrorResponse(errr,w)
	}
	resUser:=  output{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	w.WriteHeader(http.StatusOK)
	dat,err:= json.Marshal(resUser)
	if err!=nil {
		ServerErrorResponse(w)
	}
	w.Write(dat)
}