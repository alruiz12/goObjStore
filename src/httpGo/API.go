package httpGo
import(
	"net/http"
	"fmt"
	"io/ioutil"
)

func PUT(w http.ResponseWriter, r *http.Request){
	buf, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(buf))
}