package elev

import (
	"fmt"
	"time"

	"./driver/elevio"
)

func MotorController(motorChannel chan MotorDirection) { // Was int before !!!
	for {
		select {
		case command := <-motorChannel:
			switch command {
			case DIR_UP:
				elevio.SetMotorDirection(elevio.MD_Up)
			case DIR_DOWN:
				elevio.SetMotorDirection(elevio.MD_Down)
			case DIR_STOP:
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
		}
	}
}

func LightController(lightChannel chan Light) {
	for {
		select {
		case command := <-lightChannel:
			switch command.LightType {
			case BUTTON_HALL_UP:
				elevio.SetButtonLamp(elevio.BT_HallUp, command.FloorNumber, command.LightOn)
			case BUTTON_HALL_DOWN:
				elevio.SetButtonLamp(elevio.BT_HallDown, command.FloorNumber, command.LightOn)
			case BUTTON_CAB:
				elevio.SetButtonLamp(elevio.BT_Cab, command.FloorNumber, command.LightOn)
			case FLOOR_INDICATOR:
				elevio.SetFloorIndicator(command.FloorNumber)
			}
		}
	}
}

func ActionController(buttonChannel chan elevio.ButtonEvent, lightChannel chan Light, requestActionChannel chan Action, sendChannel chan ElevatorOrderMessage) {
	for {
		select {
		case buttonEvent := <-buttonChannel:
			switch buttonEvent.Button {
			case (elevio.ButtonType)(BUTTON_HALL_UP):
				fmt.Println("Up button!")
				sendChannel <- ElevatorOrderMessage{
					Event:     EVENT_NEW_ORDER,
					Direction: DIR_UP,
					Floor:     buttonEvent.Floor,
					Origin:    "node",
					Sender:    "sender",
				}

			case (elevio.ButtonType)(BUTTON_HALL_DOWN):
				fmt.Println("Down button!")
				sendChannel <- ElevatorOrderMessage{
					Event:     EVENT_NEW_ORDER,
					Direction: DIR_DOWN,
					Floor:     buttonEvent.Floor,
					Origin:    "node",
					Sender:    "sender",
				}

			case (elevio.ButtonType)(BUTTON_CAB):
				fmt.Println("Cab button!")
				lightChannel <- Light{
					LightType:   BUTTON_CAB,
					FloorNumber: buttonEvent.Floor,
					LightOn:     true,
				}
				AddCabOrder(buttonEvent.Floor)
				//UpdateReservationTable(sendChannel)
				//AddReservationOrder(buttonEvent.Floor)
			}
		case requestEvent := <-requestActionChannel:
			switch requestEvent.Command {
			case ACTION_REQUEST_ORDER:
				fmt.Println("Action: Request order!")
				CheckForOrders(sendChannel)
			case ACTION_REQUEST_SPECIFIC_ORDER:
				sendChannel <- ElevatorOrderMessage{
					Event:     EVENT_ORDER_RESERVE_SPECIFIC,
					Direction: requestEvent.Direction,
					Floor:     requestEvent.Floor,
					Origin:    "node",
					Sender:    "node",
				}
			case ACTION_ORDER_DONE:
				fmt.Println("Action: Order Done")
				sendChannel <- ElevatorOrderMessage{
					Event:  EVENT_ORDER_DONE,
					Floor:  requestEvent.Floor,
					Origin: "node",
					Sender: "node",
				}
			case ACTION_RESET_ALL_LIGHTS:
				for i := 0; i < (MAX_FLOOR_NUMBER - 1); i++ {
					lightChannel <- Light{
						LightType:   BUTTON_HALL_UP,
						FloorNumber: i,
						LightOn:     false,
					}
				}
				for i := 0; i < MAX_FLOOR_NUMBER-1; i++ {
					lightChannel <- Light{
						LightType:   BUTTON_HALL_DOWN,
						FloorNumber: i,
						LightOn:     false,
					}
				}
				for i := 0; i < MAX_FLOOR_NUMBER; i++ {
					lightChannel <- Light{
						LightType:   BUTTON_CAB,
						FloorNumber: i,
						LightOn:     false,
					}
				}
			}
		}
	}
}

func DoorController(doorChannel chan bool) {
	for {
		select {
		case openDoor := <-doorChannel:
			if openDoor {
				openDoorAction()
			} else {
				elevio.SetDoorOpenLamp(false)
			}
		}
	}
}

func openDoorAction() {
	elevio.SetDoorOpenLamp(true)
	time.Sleep(2 * time.Second)
	elevio.SetDoorOpenLamp(false)
}
