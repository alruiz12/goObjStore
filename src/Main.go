package main

import (
	"net/http"
	"log"
	"simpleBT/src/conf"
	"os"
	"io/ioutil"
	"strconv"
)


func main() {
	/*
	Initializing (and creating if non existing) files
	*/
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


	router := conf.MyNewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}