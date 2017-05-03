package tcp


import "net"
import "fmt"
import "bufio"
import "strings"
import "time"


func PeerSend(IP string, port string) {

	go func() {
		// connect to this socket
		conn, err := net.Dial("tcp", IP+port) //ex:"127.0.0.1:8081"
		if err != nil {
			fmt.Println(err.Error())
		}

		for {
			text:="1000000"
			// send to socket
			fmt.Fprintf(conn, text + "\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("Message from tracker: " + message)
		}
	}()
	time.Sleep(1 * time.Minute)
}

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