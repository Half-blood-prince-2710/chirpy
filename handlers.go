package main

import (
	"database/sql"
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
	if cfg.envi.mode != "dev"{
		w.WriteHeader(http.StatusForbidden)
	}
	cfg.fileserverHits.Store(0)
	err:=cfg.db.DeleteAllUsers(r.Context())
	if err!=nil{
		ServerErrorResponse(w)
		return
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
	
	w.Header().Set("Content-Type","application/json")
	err:= json.NewDecoder(r.Body).Decode(&input)
	if err!=nil{
		badRequestErrorResponse(w)
		return
	}
	hash,err:=auth.HashPassword(input.Password)
	if err!=nil{
		ServerErrorResponse(w)
		return
	}
	dat:=database.CreateUserParams{
		Email: input.Email,
		HashedPassword: hash,
	}
	slog.Info("Response: ","email",input.Email)
	user,err:= cfg.db.CreateUser(r.Context(),dat)
	if err!=nil{
		ServerErrorResponse(w)
		return
	}
	output.ID = user.ID
	output.CreatedAt = user.CreatedAt
	output.UpdatedAt = user.UpdatedAt
	output.Email = user.Email
	slog.Info("Response: ","user",user)

	data, err:= json.MarshalIndent(output," ","\t")
	if err!=nil {
		ServerErrorResponse(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)


}


func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	
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
	// var success struct {
	// 	Valid bool
	// }	
	
	err := json.NewDecoder(r.Body).Decode(&input)
	if err!=nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-type","application/json")
		er.Error = "Something went wrong"
		dat,err := json.MarshalIndent(er," ","\t")
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		w.Write(dat)
		return
	}
	
	// userId := r.Context().Value("userID")
	// input.UserId = uuid.Parse(userId)
	 userID, ok := r.Context().Value("userID").(uuid.UUID)
	  if !ok {
        http.Error(w, "UserID not found in context", http.StatusInternalServerError)
        return
    }
	input.UserId = userID
	fmt.Print("\n","body",input.Body,"\n userid",input.UserId,"\n")
	if len(input.Body) >140 {
			w.WriteHeader(400)
		w.Header().Set("Content-type","application/json")
		er.Error = "Chirp is too long"
		dat,err := json.MarshalIndent(er," ","\t")
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
	slog.Info("Response: ","chirp",ch)
	chirp , err:= cfg.db.CreateChirp(r.Context(),ch)
	if err!=nil{
		dbErrorReponse(err,w)
		return
	}
	var output struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
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
		// success.Valid = true
		dat,err := json.MarshalIndent(output," ","\t")
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
		
	w.Write(dat)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps,err:= cfg.db.GetAllChirps(r.Context())
	if err!=nil{
		dbErrorReponse(err,w)
		return

	}
	

	w.WriteHeader(http.StatusOK)
	data,err := json.MarshalIndent(chirps," ","\t")
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
		slog.Error("err parsing uuid")
		ServerErrorResponse(w)
		return
	}
	chirp,err:= cfg.db.GetOneChirp(r.Context(),idx)
	slog.Info("Response: ","chirp",chirp)
	dbErrorReponse(err,w)

	w.WriteHeader(http.StatusOK)
	data,err := json.MarshalIndent(chirp," ","\t")
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
		// Expires int `json:"expires_in_seconds"`
	}
	type output struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	
	
	
	w.Header().Set("Content-Type","application/json")
	err:= json.NewDecoder(r.Body).Decode(&input)
	if err!=nil{
		badRequestErrorResponse(w)
		return
	}
	// fmt.Print("login expires 1: ",input.Expires,"\n")
	//validating expire input
	// if input.Expires ==0 || input.Expires> 60*60 {
	// 	input.Expires = 3600 
	// }

	user,err:=cfg.db.FindUserByEmail(r.Context(),input.Email)
	if err!=nil {
		unauthorizedErrorResponse(w,err.Error())
		return
	}

	err =auth.CheckPasswordHash(input.Password,user.HashedPassword)
	if err!=nil {
		unauthorizedErrorResponse(w,err.Error())
		return
	}
	// fmt.Print("expireins login handler: ", input.Expires,"\n")
	token , err:= auth.MakeJWT(user.ID,cfg.envi.jwtSecret) //,time.Duration(input.Expires * int(time.Second))
	if err!=nil {
		unauthorizedErrorResponse(w,err.Error())
		slog.Error("make jwt: ","err: ",err)
		return
	}
	refreshToken ,err := auth.MakeRefreshToken()
	if err!=nil{
		slog.Error("error making refresh token ","err",err)
		ServerErrorResponse(w)
		return
	}
	data := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour*24*60),

	}
	_,err =cfg.db.CreateRefreshToken(r.Context(),data)
	if err!=nil {
		slog.Error("creating refresh at db : ","err",err)
		dbErrorReponse(err,w)
		return
	}

	resUser:=  output{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	dat,err:= json.MarshalIndent(resUser," ","\t")
	if err!=nil {
		ServerErrorResponse(w)
		return
	}
	w.Write(dat)
}


func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		unauthorizedErrorResponse(w, "no refresh token")
		return
	}

	refreshToken, err := cfg.db.GetTokenByToken(r.Context(), token)
	if err != nil {
		slog.Error("error looking up refresh token in DB", "err", err)
		dbErrorReponse(err, w)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		slog.Warn("refresh token has expired")
		unauthorizedErrorResponse(w, "refresh token expired")
		return
	}
	if refreshToken.RevokedAt.Valid {
		slog.Error("refresh token is revoked")
		unauthorizedErrorResponse(w, "revoked refresh token")
		return
	}
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.envi.jwtSecret)
	if err != nil {
		slog.Error("error creating access token", "err", err)
		unauthorizedErrorResponse(w, "error creating access token")
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("error encoding JSON response", "err", err)
		ServerErrorResponse(w)
		return
	}
}


func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		unauthorizedErrorResponse(w, "no refresh token")
		return
	}
	data := database.UpdateRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(),Valid: true},
		UpdatedAt: time.Now(),
		Token: token,

	}
	err =cfg.db.UpdateRefreshToken(r.Context(),data)
	if err!=nil {
		dbErrorReponse(err,w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}