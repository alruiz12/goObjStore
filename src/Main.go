package main

import (
	"net/http"
	"log"
	"simpleBT/src/conf"

)


func main() {
	router := conf.MyNewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}