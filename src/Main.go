package main

import (
	"net/http"
	"log"

	"github.com/alruiz12/simpleBT/src/conf"
)


func main() {
	router := conf.MyNewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}