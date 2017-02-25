package main

import (
	"fmt"
	"net/http"
	"log"
	"conf"
	//"mux"
)


func main() {
	fmt.Println("...My start...")
	//list:= []string{"http:Torrentfile1"}
	// Initialize new tracker
	// Get Meta Info from file
	// Concurrently, go routine will be blocked in the "select" waiting for a os.Interrupt or os.Kill
	// Listens stuff and handles http
	// communicate tracker with it via HTTP from peer
	// tracker.StartTracker("address",list)
	// args := flag.Args() //returns the non-flag command-line arguments.
	//err := torrent.RunTorrents(nil, args)
	/* if err != nil {log.Fatal("Could not run torrents", args, err)}*/

	var router = conf.MyNewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("...My END...")
}