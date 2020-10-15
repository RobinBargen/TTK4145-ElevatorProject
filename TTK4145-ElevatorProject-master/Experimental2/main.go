package main

import (
	//"time"

	"fmt"

	"./elev"
	"./elev/driver/elevio"
	"./elev/network/bcast"
)

func main() {
	//var id string
	elevio.Init("localhost:9999", 4)
	//id = net.SetUpLocalIP()

	elev.TargetFloor = elev.UNDEFINED_TARGET_FLOOR
	//elev.TargetFloor = -1
	motorChannel := make(chan elev.MotorDirection)
	lightChannel := make(chan elev.Light)
	doorChannel := make(chan bool)
	floorChannel := make(chan int)
	buttonChannel := make(chan elevio.ButtonEvent)

	requestChannel := make(chan elev.Action)

	sendOrderChannel := make(chan elev.ElevatorOrderMessage)
	receiveOrderChannel := make(chan elev.ElevatorOrderMessage)

	go elev.ActionController(buttonChannel, lightChannel, requestChannel, sendOrderChannel)
	go elev.FiniteStateMachine(motorChannel, lightChannel, floorChannel, doorChannel, requestChannel)
	//go elev.CheckForOrders(sendOrderChannel)

	go elev.MotorController(motorChannel)
	go elev.LightController(lightChannel)
	go elev.DoorController(doorChannel)

	go elevio.PollFloorSensor(floorChannel)
	go elevio.PollButtons(buttonChannel)

	go bcast.Transmitter(15100, sendOrderChannel)
	go bcast.Receiver(15100, receiveOrderChannel)

	//go elev.PrintElevatorInfo()

	go func() {
		for {
			select {
			case message := <-receiveOrderChannel:
				switch message.Event {
				case elev.EVENT_NEW_ORDER:
					fmt.Println("EVENT_NEW_ORDER")
					elev.NewOrderEvent(message, sendOrderChannel)
				case elev.EVENT_ACK_NEW_ORDER:
					elev.AckNewOrderEvent(message, lightChannel)

				case elev.EVENT_ORDER_RESERVE:
					elev.OrderReserveEvent(message, sendOrderChannel)
				case elev.EVENT_ACK_ORDER_RESERVE:
					elev.AckOrderReserveEvent(message)

				case elev.EVENT_ORDER_RESERVE_SPECIFIC:
					elev.OrderReserveSpecificEvent(message, sendOrderChannel)
				case elev.EVENT_ACK_ORDER_RESERVE_SPECIFIC:
					elev.AckOrderReserveSpecificEvent(message)
				case elev.EVENT_ORDER_DONE:
					elev.OrderDoneEvent(message, sendOrderChannel)
				case elev.EVENT_ACK_ORDER_DONE:
					elev.AckOrderDoneEvent(message, lightChannel)
				default:
					// Do nothing
				}
			}
		}
	}()

	//go elev.MessageLog(receiveChannel)
	//go elev.ReceiveHandler(receiveChannel, sendChannel, lightChannel)
	//time.Sleep(1 * time.Second)
	select {}
}
