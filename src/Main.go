package main

import (
	"os"
	"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"

)

func main() {

	//var filePath2 = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile"
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml"


	var proxyAddr = "127.0.0.1:8080"
	var trackerAddr = "127.0.0.1:8070"

	// tracker sends to:
	var Peer1 = "127.0.0.1:8081"
	var Peer2 = "127.0.0.1:8082"
	var Peer3 = "127.0.0.1:8083"
	var Peer4 = "127.0.0.1:8084"
	var Peer5 = "127.0.0.1:8085"
	var Peers =[]string{Peer1, Peer2, Peer3, Peer4, Peer5}
	// last character of port is the peer's internal identifier


	httpGo.StartTracker(Peers)


	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){http.ListenAndServe(":8070", routerPeer)}()
	go func(){httpGo.Put(filePath, proxyAddr, trackerAddr, 5)
		time.Sleep(5*time.Second)
		httpGo.Get("0527cbea2805d89c6d5d6457b7f9f77c",proxyAddr, trackerAddr)

	}()
	go func(){http.ListenAndServe(":8081", routerPeer)}()
	go func(){http.ListenAndServe(":8082", routerPeer)}()
	go func(){http.ListenAndServe(":8083", routerPeer)}()
	go func(){http.ListenAndServe(":8084", routerPeer)}()
	go func(){http.ListenAndServe(":8085", routerPeer)}()

	http.ListenAndServe(":8080", routerTracker)

	}

