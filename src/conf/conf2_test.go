package conf

import (
	"net/http/httptest"

	"testing"
	"fmt"
	"strings"
	"io"
	"net/http"
)

var (
	tracker2	*httptest.Server
	reader2 io.Reader
	incomingURL2 string
)

func init(){
	router := MyNewRouter()
	tracker2=httptest.NewServer(router)
	incomingURL2=fmt.Sprintf("%s/addTorrent", tracker2.URL)
	fmt.Println(incomingURL2)


}

func TestAddTorrent_AddPeer(t *testing.T) {
	torrentJson := `{"name":"torrent1"}`
	reader2 = strings.NewReader(torrentJson)
	request, err := http.NewRequest("POST", incomingURL2, reader2)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 201 {
		t.Error("Success expected: %d", res.StatusCode)
	}
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
}
