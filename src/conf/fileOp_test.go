package conf
import (
	"testing"
	"os"
	"fmt"
	"time"

)
func TestFileOp(t *testing.T){
	SplitFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	splitSuccesful:=CheckPieces("bigFile")
	if splitSuccesful==false {
		t.Failed()
		fmt.Println("Error")
	}else{fmt.Println("Good1")}
	time.Sleep(1 * time.Second)
	SplitFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	splitSuccesful=CheckPieces("bigFileInvented")
	if splitSuccesful==true {
		t.Failed()
		fmt.Println("Error2")
	}else{fmt.Println("is false")}




}