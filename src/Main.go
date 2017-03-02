package main

import (
	"fmt"
	"net/http"
	"log"
	"simpleBT/src/conf"
	//"mux"

	"os"
	"io/ioutil"
	"strconv"
//	"os/exec"
)


func main() {
	fmt.Println("...My start...")
	f, err:=os.Create("torrentsFile")
	if err!=nil {panic (err)}
	f.Close()

	f, err=os.Create("nTorrents")
	if err!=nil {panic (err)}
	nTorrents:=0
	if err = ioutil.WriteFile("nTorrents", []byte(strconv.Itoa( nTorrents)), 0666); err != nil{
		panic(err)
	}
	f.Close()
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

	router := conf.MyNewRouter()
	fmt.Println("...My END...")

	log.Fatal(http.ListenAndServe(":8080", router))



}