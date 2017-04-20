package conf
import (
	"testing"
	"os"
	"fmt"

)
func TestFileOp(t *testing.T){
	SplitFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	splitSuccesful:=CheckPieces("bigFile")
	if splitSuccesful==false {
		t.Failed()
		fmt.Println("Error")
	}else{fmt.Println("Good1")}





}