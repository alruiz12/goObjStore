package conf

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"

	"github.com/alruiz12/simpleBT/src/vars"
	"strings"
	"errors"
	"time"
)

var quit = make(chan struct{})
/*
upLoadFile is called when a POST requests 8080/upLoadFile.
Allow peer to upload a file
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request.
 */
func upLoadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*** addFile STARTS ***")
	var file string
	if r.Method == http.MethodPost{
		f, header, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file", 404)
			return
		}
		//if (Exists("../uploadedFiles/"+header.Filename)){ }
		defer f.Close()
		fileName:=header.Filename
		destination, err := os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/uploadedFiles/"+fileName)
		if err != nil {
			http.Error(w,err.Error(), 501) //internal server error
			return
		}
		defer destination.Close()
		io.Copy(destination,f)

		body, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		//file filled  with body
		file = string(body)

	}
	w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")
	fmt.Fprintln(w,`
	<form action="/upLoadFile" method="post" enctype="multipart/form-data">
	    upload a file<br>
	    <input type="file" name="file"><br>
	    <input type="submit">
	</form>
	<br>
	<br>
	`,file)

	fmt.Println("*** addFile FINISHES ***")
}

func StartAnnouncing(interval time.Duration, stopTime time.Duration){
	ticker := time.NewTicker(interval * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				announce()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	time.AfterFunc(stopTime * time.Second, closeQuit )

}

func closeQuit(){
	close(quit)
}

func announce(){
	fmt.Println("announce")
	var reader io.Reader
	trackerURL := "http://"+vars.TrackerIP+vars.TrackerPort+"/listenAnnounce"
	peerURL := trackerURL	//Variable replication just for the sake of clarity
	jsonContent := `{"file":"torrent1","IP":"`+peerURL+`"}`
	reader = strings.NewReader(jsonContent)
	request, err := http.NewRequest("POST", trackerURL, reader)
	req, err := http.DefaultClient.Do(request)
	fmt.Println("announce answer:"+ req.Status)
	if err != nil {
		errors.New("invalid request")
	}
}
