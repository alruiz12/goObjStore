package tcp


import "net"
import "fmt"
import "bufio"
import "strings"
import "time"
import (
	"io/ioutil"
	"os"
	"strconv"
)


func PeerSend(IP string, port string, filePath string) {
	var status=0	// 0 = send new size
			// 1 = resend message
			// 2 = send content

	var allBytes []byte
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	go func() {
		// connect to this socket
		conn, err := net.Dial("tcp", IP+port) //ex:"127.0.0.1:8081"
		if err != nil {
			fmt.Println(err.Error())
		}
		text:=strconv.FormatInt(fileInfo.Size(),10)	// size
		for {
			if status == 2 {	// send content
				allBytes,err =ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				text=string(allBytes)
				fmt.Println("sending file data")
			}

			// send to socket

			if status==2{
				n, err:=fmt.Fprintf(conn, text)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("peer: bytes send: ",n)
			}else{
				n, err:=fmt.Fprintf(conn, text + "\n")
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("peer: bytes send: ",n)
			}


			if status==2 {return }
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("Message from tracker: " + message)

			if strings.Compare(message,"OK\n")==0 {
				fmt.Println("Received OK")
				status = 2
			}else{
				status = 1 // resend message
			}
		}
	}()
	time.Sleep(5 * time.Minute)
}
/*
func getNextFilePart(file string) string{

}
*/
func PeerListen(port string) {

	fmt.Println("Peer listening...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}

	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err.Error())
	}
	go func() {
		for {
			// @Todo limit buffer ('\n' not valid)
			// will listen for message to process ending in newline (\n)
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
			}
			// output message received
			fmt.Print("Message Received:", string(message))

			// @Todo if metadata, save and be ready for data
				//@Todo: else, save where appropriate
			//@Todo: address other peers
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			conn.Write([]byte(newmessage + "\n"))
		}
	}()
}