package main

import (
	"./elevio"
)

func main() {
	elevio.Init("localhost:15657", 4)

	for {
		elevio.SetMotorDirection(1)
	}

}
