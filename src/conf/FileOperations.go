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
var totalPartsNum uint64
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
	totalPartsNum= uint64(math.Ceil(float64(fileSize)/float64(fileChunk)))
	fmt.Printf("File size= %d  Bytes\n",fileSize)
	fmt.Printf("Splitting to %d pieces \n", totalPartsNum)
	var partSize int
	var partBuffer []byte
	var newFileName string
	allBytes,_:=ioutil.ReadFile(filePath)
	fmt.Println(len(allBytes))
	for i := uint64(0); i<totalPartsNum; i++{
		partSize=int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer=make([]byte,partSize)
		partBuffer=allBytes[:partSize]
		allBytes=allBytes[partSize:]
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("______________________________________n bytes read: ",len(partBuffer)," out of: ",partSize)
		fmt.Printf("first char is: %s\n", string(partBuffer[1]))
		fmt.Println(string(partBuffer))

		// write to disk
		newFileName= path.Base(filePath)+"_"+strconv.FormatUint(i,10)+"_"
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
	_, err =  os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/z"+fileName)
	newFile, err:=os.OpenFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/z"+fileName,os.O_APPEND|os.O_WRONLY,0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//defer newFile.Close()
	var inOrderCount uint64=0
	for inOrderCount<totalPartsNum {
		for _, file := range files {
			if strings.Compare(file.Name(), fileName + "_" + strconv.FormatUint(inOrderCount,10) + "_") == 0 {
				inOrderCount++
				fmt.Println("-------------------------- ****file:= " + file.Name())

				currentFile, err := os.Open(path + file.Name())
				if err != nil {
					fmt.Println(err)
					return false
				}
				defer currentFile.Close()
				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())
				fmt.Println("bytes current file: ", string(bytesCurrentFile))
				_, err= newFile.WriteString(string(bytesCurrentFile))
				if err != nil {
					panic(err)
				}
				//ioutil.WriteFile(path+"z"+fileName, bytesCurrentFile, os.ModeAppend)
				//newFile.Write(bytesCurrentFile)
			}
		}
	}

	return true
}
