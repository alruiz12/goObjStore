package main
import()
import (
	"net/http/httptest"
	"github.com/alruiz12/simpleBT/src/conf"
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
	router := conf.MyNewRouter()
	tracker=httptest.NewServer(router)
	incomingURL=fmt.Sprintf("%s/addTorrent", tracker.URL)
	fmt.Println(incomingURL)
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