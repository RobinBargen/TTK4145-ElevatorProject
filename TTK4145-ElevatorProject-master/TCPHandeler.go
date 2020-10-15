import (
  "fmt"
  "net"
)

func sendCommand(message string){
  connection, error := net.Dial("tcp", "129.241.187.255:34933")
  if(error != nil){
    fmt.Println("Error in sendCommand(message string)!")
  }
  connection.Write([]byte("Hei"))
  buffer := make([]byte, 1024)
  ln, readError := connection.Read(buffer)
  if(error != nil){
    fmt.Println("Error in sendCommand(message string)!")
  }
  fmt.Println(string(buffer[0:ln])))
}
