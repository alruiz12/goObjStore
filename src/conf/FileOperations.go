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
var totalPartsNum uint64
const fileChunk = 1*(1<<10) // 1 KB
var mainFileHash string


/*
SplitFile opens the file we want to split, computes the original hash,
and saves each chunk as a new file following the fashion "bigFile_0_"
@param path to the file we want to split
*/
func SplitFile(filePath string){
	// Open the file we want to split
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Compute original hash
	hash := md5.New()
	_, err = io.Copy(hash,file)
	if err != nil {
		fmt.Println(err)
		return
	}
	mainFileHash=hex.EncodeToString(hash.Sum(nil))
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

		// create new file
		newFileName= path.Base(filePath)+"_"+strconv.FormatUint(i,10)+"_"
		_, err =  os.Create(newFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		// write / save buffer to file
		ioutil.WriteFile(newFileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to: ",newFileName)

	}


}

/*
md5sum opens the file we want to compute the hash and computes it
@param path to the file we want to split
returns the computed hash
*/
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
	return mainFileHash
}

/*
CheckPieces walks through the subfiles directory, creates a new file to be filled out with the content of each subfile,
and compares the new hash with the original one.
@param path to the file we want to split
Returns true if both hash are identic and false if not
*/
func CheckPieces(fileName string) bool{

	// Subfiles directory
	path:=os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Create new file
	_, err =  os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/z"+fileName)
	newFile, err:=os.OpenFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/z"+fileName,os.O_APPEND|os.O_WRONLY,0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer newFile.Close()

	// Trying to fill out the new file using subfiles (in order)
	var inOrderCount uint64=0
	for inOrderCount<totalPartsNum {
		for _, file := range files {
			if strings.Compare(file.Name(), fileName + "_" + strconv.FormatUint(inOrderCount,10) + "_") == 0 {
				inOrderCount++

				currentFile, err := os.Open(path + file.Name())
				if err != nil {
					fmt.Println(err)
					return false
				}
				defer currentFile.Close()

				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())

				_, err= newFile.WriteString(string(bytesCurrentFile))
				if err != nil {
					fmt.Println(err)
					return false
				}

			}
		}
	}

	// Compute and compare new hash
	newHash:=md5sum(path+"/z"+fileName)

	if strings.Compare(mainFileHash,newHash)!=0 {return false}
	return true
}
