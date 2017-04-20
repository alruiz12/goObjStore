package conf
import (
	"testing"
	"os"
	"fmt"

)
func TestFileOp2(t *testing.T){

	SplitFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	splitSuccesful:=CheckPieces("bigFileInvented")
	if splitSuccesful==true {
		t.Failed()
		fmt.Println("Error2")
	}else{fmt.Println("is false")}




}