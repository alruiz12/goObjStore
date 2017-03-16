package conf

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"
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
		destination, err := os.Create("../uploadedFiles/"+fileName)
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
	    <input type="file" name="userFile"><br>
	    <input type="submit">
	</form>
	<br>
	<br>
	`,file)

	fmt.Println("*** addFile FINISHES ***")
}

/* Next Commit:
func Exists (name string) bool {
	if _, err:= os.Stat(name); err!= nil {
		if os.IsNotExist(err){
			return false
		}
	}
	return true
}

func announce(){
	var reader io.Reader
	var outgoingURL string
	jsonContent := `{"name":"torrent1"}`
	reader = strings.NewReader(jsonContent)
	request, err := http.NewRequest("POST", outgoingURL, reader)
	_, err = http.DefaultClient.Do(request)
	if err != nil {
		errors.New("invalid request")
	}
}
*/