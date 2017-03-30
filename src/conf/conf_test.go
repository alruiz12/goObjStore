package conf

import (
	"net/http/httptest"

	"testing"
	"fmt"
	"strings"
	"io"
	"net/http"
	"mime/multipart"
	"bytes"
	"strconv"
	"os"

)

var (
	tracker2	*httptest.Server
	reader2 io.Reader
	incomingURL2 string
)

func init(){



}

func TestRepo(t *testing.T) {
	router := MyNewRouter()
	tracker2=httptest.NewServer(router)


	incomingURL2=fmt.Sprintf("%s/getTorrentsList", tracker2.URL)
	fmt.Println(incomingURL2)
	torrentJson := `{"name":"torrent1"}`
	reader2 = strings.NewReader(torrentJson)
	request, err := http.NewRequest("GET", incomingURL2, reader2)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Error("Failure expected: %d", res.StatusCode)
	}


	incomingURL2=fmt.Sprintf("%s/addTorrent", tracker2.URL)
	fmt.Println(incomingURL2)
	torrentJson = `{"name":"torrent1"}`
	reader2 = strings.NewReader(torrentJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 201 {
		t.Error("Success expected: %d", res.StatusCode)
	}

	incomingURL2=fmt.Sprintf("%s/addPeer", tracker2.URL)
	fmt.Println(incomingURL2)
	peerJson := `{"peerIP":"192.168.1.3","torrentName":"torrent1"}`
	reader2 = strings.NewReader(peerJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 201 {
		t.Error("Success expected: %d", res.StatusCode)
	}




	incomingURL2=fmt.Sprintf("%s/getIPs", tracker2.URL)
	fmt.Println(incomingURL2)
	getIPJson := `{"name":"torrent1"}`
	reader2 = strings.NewReader(getIPJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200{
		t.Error("Success expected: %d", res.StatusCode)
	}


	incomingURL2=fmt.Sprintf("%s/addPeer", tracker2.URL)
	fmt.Println(incomingURL2)
	peerJson = `{"peerIP":"192.168.1.3","torrentName":"nonExistingTorrentOnPurpose"}`
	reader2 = strings.NewReader(peerJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Error("Failure expected: %d", res.StatusCode)
	}


	incomingURL2=fmt.Sprintf("%s/getIPs", tracker2.URL)
	fmt.Println(incomingURL2)
	getIPJson = `{"name":"nonExistingTorrentOnPurpose"}`
	reader2 = strings.NewReader(getIPJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404{
		t.Error("Failure expected: %d", res.StatusCode)
	}


	incomingURL2=fmt.Sprintf("%s/addTorrent", tracker2.URL)
	fmt.Println(incomingURL2)
	reader2 = nil
	request, err = http.NewRequest("GET", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200{
		t.Error("Success expected: %d", res.StatusCode)
	}

	incomingURL2=fmt.Sprintf("%s/addTorrent", tracker2.URL)
	fmt.Println(incomingURL2)
	torrentJson = `{"name":"torrent1"}`
	reader2 = strings.NewReader(torrentJson)
	request, err = http.NewRequest("POST", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 208 {
		t.Error("Failure expected: %d", res.StatusCode)
	}


	incomingURL2=fmt.Sprintf("%s/getTorrentsList", tracker2.URL)
	fmt.Println(incomingURL2)
	torrentJson = `{"name":"torrent1"}`
	reader2 = strings.NewReader(torrentJson)
	request, err = http.NewRequest("GET", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Success expected: %d", res.StatusCode)
	}







	file, err := os.Open("uploadedFile")
	if err != nil {
		fmt.Println("Opening file")
		t.Error(err)
	}
	defer file.Close()
	incomingURL2=fmt.Sprintf("%s/upLoadFile", tracker2.URL)
	fmt.Println(incomingURL2)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "uploadedFile")
	if err != nil {
		fmt.Println("creating Form file")
		t.Error(err)
	}
	_, err = io.Copy(part, file)
	err=writer.Close()
	if err != nil {
		t.Error(err)
	}
	request, err = http.NewRequest("POST", incomingURL2, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Success expected: %d", res.StatusCode)
	}











	var port string =tracker2.URL[len(tracker2.URL)-5:]
	var IPaddr = tracker2.URL[7:len(tracker2.URL)-6]
	sendFile("torrent1",IPaddr,port)


	incomingURL2=fmt.Sprintf("%s/p2pRequest", tracker2.URL)
	fmt.Println(incomingURL2)
	status:=2
	torrentJson = `{"file":"torrent1","IP":"`+IPaddr+`","status":"`+strconv.Itoa(status)+`","port":"`+port+`"}`
	reader2 = strings.NewReader(torrentJson)
	request, err = http.NewRequest("GET", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Success expected: %d", res.StatusCode)
	}



	incomingURL2=fmt.Sprintf("%s/p2pRequest", tracker2.URL)
	fmt.Println(incomingURL2)
	status=2
	torrentJson = `{"file":"NON_EXISTING_torrent","IP":"`+IPaddr+`","status":"`+strconv.Itoa(status)+`","port":"`+port+`"}`
	reader2 = strings.NewReader(torrentJson)
	request, err = http.NewRequest("GET", incomingURL2, reader2)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Success expected: %d", res.StatusCode)
	}


}
