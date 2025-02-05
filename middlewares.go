package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/half-blood-prince-2710/chirpy/internal/auth"
)


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

// func middlewareLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }


func (cfg *apiConfig) authenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Print("entering auth middle\n")
		token , err:=auth.GetBearerToken(r.Header)
		if err!=nil {
			slog.Error("authenticate middleware : " ,"err",err)
			unauthorizedErrorResponse(w,err.Error())
			slog.Error("Error: get bearer token ","err",err)
			return 
		}
		// fmt.Print("entering auth middle\ntoken:",token,"\n")
		id,err:=auth.ValidateJWT(token,cfg.envi.jwtSecret)
		if err!=nil {
			slog.Error("authenticate middleware : " ,"err",err)
			
			unauthorizedErrorResponse(w, err.Error())
			return 
		}
		// fmt.Print("entering auth middle\nid: ",id,"\n")
		ctx := context.WithValue(r.Context(), "userID", id)
		req := r.WithContext(ctx)
		next.ServeHTTP(w,req)
	}) 

		
}
