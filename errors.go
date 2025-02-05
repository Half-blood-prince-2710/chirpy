package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorRes struct {
		Msg string `json:"error"`
	}

// type envelope []map[string]any

func ServerErrorResponse(w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
	
}
func badRequestErrorResponse(w http.ResponseWriter){

	w.WriteHeader(http.StatusBadRequest)
	errRes := errorRes{Msg:"Bad Request"}
	dat,err:=json.Marshal(errRes)
	if err!=nil{
		ServerErrorResponse(w)
	}
	w.Write(dat)
}
func dbErrorReponse(err error,w http.ResponseWriter) {
	
		slog.Error("Error in db: ","err: ",err)
		w.WriteHeader(404)
		ServerErrorResponse(w)
	
}

func unauthorizedErrorResponse(w http.ResponseWriter,msg string) {
	
	errr :=  errorRes{Msg:msg}
	dat,err:=json.Marshal(errr)
	if err!=nil {
		ServerErrorResponse(w)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
		w.Write(dat)
}