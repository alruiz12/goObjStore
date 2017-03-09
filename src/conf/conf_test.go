package conf
import()
import (
	"net/http/httptest"

	"testing"
	"fmt"
	"strings"
	"io"
	"net/http"
)

var (
	tracker	*httptest.Server
	reader io.Reader
	incomingURL string
)

func init(){
	router := MyNewRouter()
	tracker=httptest.NewServer(router)

	incomingURL=fmt.Sprintf("%s/addTorrent", tracker.URL)
	fmt.Println(incomingURL)
	fmt.Println("yeseeeeeeeee")/*
	incomingURL=fmt.Sprintf("%s/addPeer", tracker.URL)
	fmt.Println(incomingURL)
	fmt.Println("yeseeeeeeeee")*/
}

func TestAddTorrent(t *testing.T){
	torrentJson := `{"name":"torrent1"}`
	reader = strings.NewReader(torrentJson)
	request, err := http.NewRequest("POST",incomingURL, reader)
	res, err := http.DefaultClient.Do(request)
	if err!=nil {t.Error(err)}
	if res.StatusCode != 201 {
		t.Error("Success expected: %d",res.StatusCode)
	}
}
func TestAddPeer(t *testing.T){
	peerJson := `{"peerIP":"192.168.1.3","torrentName":"torrent1"}`
	reader = strings.NewReader(peerJson)
	request, err := http.NewRequest("POST",incomingURL, reader)
	res, err := http.DefaultClient.Do(request)
	if err!=nil {t.Error(err)}
	if res.StatusCode != 201 {
		t.Error("Success expected: %d",res.StatusCode)
	}
}