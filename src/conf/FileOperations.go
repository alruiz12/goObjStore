package conf

import (
	"os"
	"fmt"
	"math"
	"crypto/md5"
	"io"
	"encoding/hex"
	"path"
	"strconv"
	"io/ioutil"
	"strings"
)
var partSize int
var partBuffer []byte
const fileChunk = 1*(1<<10) // 1 KB

func SplitFile(filePath string){
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash,file)
	if err != nil {
		fmt.Println(err)
		return
	}
	mainFileHash:=hex.EncodeToString(hash.Sum(nil))
	fmt.Println(mainFileHash)

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	// calculate total number of parts the file will be chunked into
	totalPartsNum:= uint64(math.Ceil(float64(fileSize)/float64(fileChunk)))
	fmt.Printf("File size= %d  Bytes\n",fileSize)
	fmt.Printf("Splitting to %d pieces \n", totalPartsNum)
	var partSize int
	var partBuffer []byte
	var newFileName string
	for i := uint64(0); i<totalPartsNum; i++{
		partSize=int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer=make([]byte,partSize)
		file.Read(partBuffer)

		// write to disk
		newFileName= path.Base(filePath)+"_"+strconv.FormatUint(i,10)
		_, err =  os.Create(newFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		// write / save buffer to disk
		ioutil.WriteFile(newFileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to: ",newFileName)

	}


}

func md5sum(filePath string) string{
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash,file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	mainFileHash:=hex.EncodeToString(hash.Sum(nil))
	fmt.Println(mainFileHash)
	return mainFileHash
}

func CheckPieces(fileName string) bool{


	path:=os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err =  os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/"+fileName)
	newFile, err:=os.Open(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/"+fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer newFile.Close()
	for _, file  := range files {
		if strings.Contains(file.Name(), fileName){
			fmt.Println("file:= "+file.Name())

			currentFile, err:=os.Open(path+file.Name())
			if err != nil {
				fmt.Println(err)
				return false
			}
			defer newFile.Close()

			partSize=int(fileChunk)
			partBuffer=make([]byte,partSize)
			currentFile.Read(partBuffer)

			//contentBytes,_:=ioutil.ReadFile(path+file.Name())
			newFile.Write(partBuffer)
			fmt.Println(hex.EncodeToString(partBuffer))
		}
	}

	return true
}
