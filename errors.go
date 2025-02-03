package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errror struct {
		msg string `json="error"`
	}

type envelope []map[string]any

func ServerErrorResponse(w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
func badRequestErrorResponse(w http.ResponseWriter,  status int , errr errror){
	w.WriteHeader(http.StatusBadRequest)
	errr.msg = `Bad Request`
	dat,err:=json.Marshal(errr)
	if err!=nil{
		ServerErrorResponse(w)
	}
	w.Write(dat)
}
func dbErrorReponse(err error,w http.ResponseWriter) {
	if err!=nil{
		slog.Error("err: ",err)
		ServerErrorResponse(w)
		return
	}
}