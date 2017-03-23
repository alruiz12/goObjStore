package main

import (
	"net/http"
	"log"
	"github.com/alruiz12/simpleBT/src/conf"
	"github.com/alruiz12/simpleBT/src/vars"
)

func main() {

	router := conf.MyNewRouter()
	conf.StartAnnouncing(2,9)
	conf.CheckInactivePeers(5)
	log.Fatal(http.ListenAndServe(vars.TrackerPort, router))
}
