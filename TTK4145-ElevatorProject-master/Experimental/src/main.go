package main

import (
	//"time"

	"./elev"
	"./elev/driver/elevio"
	"./net"
	"./net/network/bcast"
)

func main() {
	//var id string
	elevio.Init("localhost:9999", 4)
	//id = net.SetUpLocalIP()

	elev.TargetFloor = elev.UNDEFINED_TARGET_FLOOR

	motorChannel := make(chan int)
	lightChannel := make(chan elev.Light)
	doorChannel := make(chan bool)

	floorChannel := make(chan int)
	buttonChannel := make(chan elevio.ButtonEvent)

	requestChannel := make(chan elev.Action) // Used by elevator to issue wanted action

	// Used to send and receive messages over the network
	sendChannel := make(chan elev.Action)
	receiveChannel := make(chan elev.Action)

	// Used to block messages sent to elevator while unavailable/busy
	//peerUpdateChannel := make(chan peers.PeerUpdate)
	//peerTxEnableChannel := make(chan bool)

	go elev.FiniteStateMachine(motorChannel, lightChannel, floorChannel, doorChannel)

	go elev.MotorController(motorChannel)
	go elev.LightController(lightChannel)
	go elev.DoorController(doorChannel)
	go elev.ActionController(buttonChannel, requestChannel, sendChannel)

	go elevio.PollFloorSensor(floorChannel)
	go elevio.PollButtons(buttonChannel)

	go bcast.Transmitter(15647, sendChannel)
	go bcast.Receiver(15647, receiveChannel)
	//go peers.Transmitter(15647, id, peerTxEnableChannel)
	//go peers.Receiver(15647, peerUpdateChannel)

	go net.MessageLog(receiveChannel)
	go net.ReceiveHandler(receiveChannel)

	//time.Sleep(1 * time.Second)
	select {}
}
