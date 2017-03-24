package main

import (
	"net/http"
	"log"
	"github.com/alruiz12/simpleBT/src/conf"
	"github.com/alruiz12/simpleBT/src/vars"
	"time"
)

func main() {

	router := conf.MyNewRouter()

	// Seeder
	go func() {
		var quit = make(chan int)
		conf.StartAnnouncing(2,9,"192.168.0.1","torrent1",quit)
		time.AfterFunc(9 * time.Second, func(){close(quit)})
		time.Sleep(5*time.Second)
	}()

	// Peer1
	go func() {
		var quit = make(chan int)
		conf.StartAnnouncing(2,15,"192.168.0.2","torrent1", quit)
		time.AfterFunc(15 * time.Second, func(){close(quit)})

	}()

	// Peer2
	go func(){
		var quit = make(chan int)
		conf.StartAnnouncing(2,21,"192.168.0.3","torrent1",quit )
		time.AfterFunc(21 * time.Second, func(){close(quit)})
	}()
	go func() {
		conf.CheckInactivePeers(5)
	}()
	log.Fatal(http.ListenAndServe(vars.TrackerPort, router))
}
