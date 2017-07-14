package httpGo
import(
	"net/http"
	"fmt"
	"os"
	"io"
	"crypto/md5"
	"encoding/hex"
	"github.com/alruiz12/simpleBT/src/conf"
	"time"
)

func PUT(w http.ResponseWriter, r *http.Request){
	var startPUT time.Time
	startPUT=time.Now()
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
	putOK := make(chan bool)
	go Put(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/"+name, conf.TrackerAddr, conf.NumNodes, putOK)
	success:= <-putOK
	if success == true {
		fmt.Println("put success ",time.Since(startPUT))
		w.WriteHeader(http.StatusCreated )
	}else{	fmt.Println("put fail")
		w.WriteHeader(http.StatusUnprocessableEntity)}
	currentKey:=md5sum(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/"+name)
	fmt.Println(CheckPieces(currentKey ,"NEW.xml",conf.FilePath, conf.NumNodes))


}


func md5String(str string) string{
	hasher:=md5.New()
	_, err:= hasher.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}