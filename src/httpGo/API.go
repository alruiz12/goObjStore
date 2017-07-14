package httpGo
import(
	"net/http"
	"fmt"
	"os"
	"io"
	"crypto/md5"
	"encoding/hex"
	"github.com/alruiz12/simpleBT/src/conf"
)

func PUT(w http.ResponseWriter, r *http.Request){
	name:=md5String(r.URL.Path)
	file, err := os.Create(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/"+name )
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(file, r.Body )
	if err != nil {
		fmt.Println(err)
	}
	file.Close()
	putDone := make(chan int)
	go Put(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/"+name, conf.TrackerAddr, conf.NumNodes, putDone)
	<-putDone

}


func md5String(str string) string{
	hasher:=md5.New()
	_, err:= hasher.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}