package conf

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"
	"bytes"
	"github.com/alruiz12/simpleBT/src/vars"
	"strings"
	"errors"
	"time"
	"encoding/json"
	"mime/multipart"
)

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

func StartAnnouncing(interval time.Duration, stopTime time.Duration,IP string,torrentName string, quit chan int){
	ticker := time.NewTicker(interval * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				announce(IP, torrentName)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()


	fmt.Println("							EXITING 22 "+IP)

}


func announce(IP string,torrentName string){
	fmt.Println("announce: ",IP)
	var reader io.Reader
	trackerURL := "http://"+vars.TrackerIP+vars.TrackerPort+"/listenAnnounce"
	//peerURL := trackerURL	//Variable replication just for the sake of clarity
	jsonContent := `{"file":"`+torrentName+`","IP":"`+IP+`"}`
	reader = strings.NewReader(jsonContent)
	request, err := http.NewRequest("POST", trackerURL, reader)
	req, err := http.DefaultClient.Do(request)
	//fmt.Println("announce answer:"+ req.Status)

	var swarmSlice []string

	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := req.Body.Close(); err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &swarmSlice);
	if err != nil {
		errors.New("error decoding swarmSlice")
	}
	fmt.Println("									SLICE: ",swarmSlice)
	var peerURL string
	jsonContent = `{"file":"`+torrentName+`","IP":"`+IP+`"}`
	reader = strings.NewReader(jsonContent)
	for _, peerIP:=range swarmSlice{
		go func(peerURL string, request *http.Request, err error, req *http.Response) {
			peerURL = "http://" + peerIP + vars.TrackerPort + "/p2pRequest"
			request, err = http.NewRequest("GET", peerURL, reader)
			req, err = http.DefaultClient.Do(request)
			//fmt.Println("p2p	p2p	p2p	p2p	p2p	p2p	p2p:" + req.Status + " by " + peerIP)
		}(peerURL,request,err,req)
	}
}


func p2pRequest(w http.ResponseWriter, r *http.Request){
	var announcement vars.Announcement
	fmt.Println("...p2pRequest starts ...")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &announcement); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println("							"+vars.IP+" was asked by "+announcement.IP)
	if strings.Compare(vars.IP, "10.0.0.11" ){ sendFile(announcement.File, announcement.IP)}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(announcement); err != nil {
		panic(err)
	}
	fmt.Println("...p2pRequest finishes ...")
}


func sendFile(fileName string, IP string){

	file, err := os.Open("src/conf/"+fileName)
	if err != nil {
		fmt.Println("Opening file")
		log.Println(err)
	}
	defer file.Close()

	//destinationURL:=fmt.Sprintf("%s/upLoadFile", IP)
	destinationURL:="http://"+IP+":8080/upLoadFile"
	//destinationURL="http://"+destinationURL
	fmt.Println(destinationURL)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("creating Form file")
		log.Println(err)
	}
	_, err = io.Copy(part, file)
	err=writer.Close()
	if err != nil {
		log.Println(err)
	}
	request, err := http.NewRequest("POST", destinationURL, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 200 {
		log.Println("Success expected: %d", res.StatusCode)
	}


}