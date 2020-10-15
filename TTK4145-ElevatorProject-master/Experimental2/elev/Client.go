package elev

/*
import (
	"fmt"
	"os"

	"./network/localip"
)

func SetUpLocalIP() string {
	var id string
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}

func MessageLog(receiveChannel chan Action) {
	for {
		select {
		case command := <-receiveChannel:
			fmt.Println(command.Parameters)
		}
	}
}
*/
