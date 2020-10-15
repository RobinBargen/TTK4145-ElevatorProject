package elev

import (
	"fmt"
	"strconv"
	"time"

	"./driver/elevio"
)

/*-----------------------------------------------------
Function:		MotorController
Arguments:	motorChannel chan int
Affected:		None
Start the motors in the specified direction,
given by the argument.
-----------------------------------------------------*/
func MotorController(motorChannel chan int) {
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

/*-----------------------------------------------------
Function:		LightController
Arguments:	lightChannel chan Light
Affected:		None
Turns on or off the lamp specified by the input argument.
-----------------------------------------------------*/
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

/*-----------------------------------------------------
Function:		ActionController
Arguments:	buttonChannel chan elevio.ButtonEvent,
						requestChannel chan Action, sendChannel chan Action
Affected:		sendChannel
Sets up the content of the sendChannel to be sent to the master - node.
-----------------------------------------------------*/
func ActionController(buttonChannel chan elevio.ButtonEvent, requestChannel chan Action, sendChannel chan Action) {
	for {
		select {
		case buttonEvent := <-buttonChannel:
			switch buttonEvent.Button {
			case (elevio.ButtonType)(BUTTON_HALL_UP):
				fmt.Println("Up button")
				parameters := "UP " + strconv.Itoa(buttonEvent.Floor)
				sendChannel <- Action{Command: COMMAND_ORDER, Parameters: parameters}
			case (elevio.ButtonType)(BUTTON_HALL_DOWN):
				fmt.Println("Down button")
				parameters := "DOWN " + strconv.Itoa(buttonEvent.Floor)
				sendChannel <- Action{Command: COMMAND_ORDER, Parameters: parameters}
			}
		case requestEvent := <-requestChannel:
			switch requestEvent.Command {
			case COMMAND_REQUEST_ORDER:
				sendChannel <- Action{Command: requestEvent.Command, Parameters: requestEvent.Parameters}
			}
		}
	}
}

/*-----------------------------------------------------
Function:		DoorController
Arguments:	doorChannel chan bool
Affected:		sendChannel
Opens the door for 2 seconds, with the help of openDoorAction(),
or closes the door. Which of the two actions performed, is dependent on
the input argument.
-----------------------------------------------------*/
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
