package main

import (
	"net/http"
	"log"
	"github.com/alruiz12/simpleBT/src/conf"
	"github.com/alruiz12/simpleBT/src/vars"
)

func main() {
	router := conf.MyNewRouter()
	log.Fatal(http.ListenAndServe(vars.TrackerPort, router))

}
