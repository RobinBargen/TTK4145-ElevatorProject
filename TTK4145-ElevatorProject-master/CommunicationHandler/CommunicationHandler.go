package CommunicationHandler

import (
	"fmt"
	"net"
)

func sendCommand(message string) {

	// Setup TCP connection
	// Raddr is our adress
	var localHostAddress string = ""
	var serverAddress string = ""
	raddr, error = net.ResolveTCPAddr("tcp", localHostAddress)
	if error != nil {
		fmt.Println("Error: sendCommand(message string) -- net.ResolveTCPAddr")
	}
	connection, error := net.DialTCP("tcp", nil, raddr)
	if error != nil {
		fmt.Println("Error: sendCommand(message string) --net.DialTCP")
	}
	// TCP - Listen
	ladder, error := net.ResolveTCPAddr("tcp", serverAddress)
	receive, error := net.ListenTCP("tcp", ladder)
	if error != nil {
		fmt.Println("Error: sendCommand(message string) --net.ListenTCP")
	}
	connection.Write([]byte("Connect to: "))

	// TCP Accept connection:
	receiveConnection, error := receive.AcceptTCP()

	buffer := make([]byte, 1024)
	messageLength, error := receiveConnection.Read(buffer)
	if error != nil {
		fmt.Println("Error: sendCommand(message string) --receiveConnection.Read")
	}
	receiveConnection.Write([]byte(message))
}
