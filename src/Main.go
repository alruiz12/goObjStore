package main

import (
	"time"
	"github.com/alruiz12/goObjSotre/src/httpGo"
	"net/http"
	"github.com/alruiz12/goObjStore/src/conf"

)
func main() {

	router := httpGo.MyNewRouter()

	go func(){http.ListenAndServe(conf.Peer2a, router)}()
	go func(){http.ListenAndServe(conf.Peer2b, router)}()
	http.ListenAndServe(conf.Peer2c, router)
	time.Sleep(1*time.Hour)

}

