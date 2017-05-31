package main

import (
	"os"
	//"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"
)

func main() {
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile"

	var trackerAddr = ":8080"

	// tracker sends to:
	var peer1 = "127.0.0.1:8081"
	var peer2 = "127.0.0.1:8082"
	var peer3 = "127.0.0.1:8083"
	peers :=[]string{peer1, peer2, peer3}
	// last character of port is the peer's internal identifier

	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){httpGo.TrackerDivideLoad(filePath, trackerAddr ,peers)}()
	go func(){http.ListenAndServe(":8081", routerPeer)}()
	http.ListenAndServe(":8080", routerTracker)
	}

