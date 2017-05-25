package tcp
import (
	"net"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
	"math"
//	"encoding/json"
)
//const fileChunk = 1*(1<<10) // 1 KB
//const fileChunk = 8*(1<<20) // 8 MB


func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//const fileChunk = 1*(1<<10) // 1 KB
/*
type peer struct {
	addr	*net.TCPAddr
	conn	*net.TCPConn
}*/
func TrackerDivideLoad(IP string, ports []int, filePath string) {
	var status=0
	// 0 = send new size
	// 2 = send content

	var n int

	go func() {
		var err error
		// create list of (port, connection)
		peers := make([]peer,len(ports))
		for index, port := range ports{
			peers[index].addr, err =net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port) )
			if err != nil {
				fmt.Println("Tracker: ResolveTCPADDr ERROR: ", err.Error())
			}
			// connect to this socket
			peers[index].conn,err =net.DialTCP("tcp4", nil ,peers[index].addr)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(" DialTCP error	"+err.Error())
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		fileInfo, _ := file.Stat()
		text:=strconv.FormatInt(fileInfo.Size(),10)	// size
		fmt.Println("SIZE: ",text)
		size,_:=strconv.Atoi(text)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//totalPartsNum:= int(math.Ceil(float64(size)/float64(fileChunk)))
		var partSize int
		var partBuffer []byte
		currentPart := 0
		var currentPeer peer
		var currentNum int=0

		/*
		var hashList []string
		//------------------------------------------------------------------------------------------------------------------------------
		for currentPart<totalPartsNum{
			// read and store hash of content in list
			partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,partSize)
			_,err = file.Read(partBuffer)
			if err != nil {
				fmt.Println("Tracker: error reading ", err.Error())
			}
			hash:=getMd5sum(partBuffer)
			hashList=append(hashList, hash)
			currentPart++
			fmt.Println("tracker current hash: ",hash)
		}

		m, err := json.Marshal(hashList)
		fmt.Println("tracker len hashlist ",len(hashList))
		if err != nil {
			fmt.Println("Error marshalling list ", err.Error())
		}

		i:=0
		for i<3 {
			_, err = fmt.Fprintf(peers[i].conn, string(m)+string('\n') )
			if err != nil {
				fmt.Println("Error sending map ", err.Error())
			}
			i++
		}
		fmt.Println("tracker hashList sent!!")

		currentPart=0
		file.Close()
		// re open file to start reading from the beginning
		file, err = os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		*/
		defer file.Close()
		// -----------------------------------------------------------------------------------------------------------------------------
		for currentPart<totalPartsNum{
			if status == 2 {	// send content
				currentPeer=peers[currentNum]
				partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
				partBuffer=make([]byte,partSize)
				_,err = file.Read(partBuffer)
				content:=string(partBuffer)
				// send to sockets
				n, err = fmt.Fprintf(currentPeer.conn, content)
				if err != nil {
					fmt.Println("ERROR sendig content !!!!!!!! ",err.Error())
				}
				fmt.Println("Tracker: n= "+ strconv.Itoa( n ) +" currentPart: "+ strconv.Itoa(currentPart) +" current peer"+ strconv.Itoa( currentPeer.addr.Port) +"|"+ content[:10])

				currentPart++
				currentNum=(currentNum+1)%3

			}else{	// status == 0 -> send size (already stored in 'text')
				// no need to update 'text'
				for _, peer := range peers {
					_, err = fmt.Fprintf(peer.conn, text + string('|'))
					if err != nil {
						fmt.Println("Error sending size")
						fmt.Println(err.Error())
					}
				}
				status=2
			}

		} // <-- for


	}() // <-- go func
	fmt.Println("tracker finishing ....................................")
	time.Sleep(15 * time.Minute)
}


/*
computes the md5 hash for the string given
@param path to the file we want to split
returns the computed hash
*/
func getMd5sum(data []byte) string{
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))

}