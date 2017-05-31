package httpGo
import(
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"
	"encoding/json"
)
func PeerListen(w http.ResponseWriter, r *http.Request){
	// Get peer ID

	// Create folder for chunks


	// Listen to tracker
	fmt.Println("*** addFile STARTS ***")
	var file string
	if r.Method == http.MethodPost{



		var announcement string
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
		fmt.Println(announcement)

		return
		//Todo pick 1 method
		f, header, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file", 404)
			return
		}
		//if (Exists("../uploadedFiles/"+header.Filename)){ }
		defer f.Close()
		fileName:=header.Filename
		destination, err := os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/httpReceived/"+fileName)
		if err != nil {
			http.Error(w,err.Error(), 501) //internal server error
			return
		}
		defer destination.Close()
		io.Copy(destination,f)

		body, err = ioutil.ReadAll(f)
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
	<form action="/upLoadFile" method="post" enctype="application/json">
	    upload a file<br>
	    <input type="file" name="file"><br>
	    <input type="submit">
	</form>
	<br>
	<br>
	`,file)

	fmt.Println("*** addFile FINISHES ***")

		// Save chunk

		// Send chunk to peers

	// Listen to other peers
}
