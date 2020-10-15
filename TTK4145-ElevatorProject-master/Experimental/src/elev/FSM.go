package elev

import (
	"fmt"
	"time"
)

var isInitialized bool = false
var isOrderServed bool = false

var state int
var previousState int

func initState(motorChannel chan int, lightChannel chan Light, floorChannel chan int) {
	fmt.Println("Init!!!")
	if isInitialized {
		state = IDLE
	}

	var elevatorIsApproaching bool = false
	for {
		floor := ReadFloorSensor(floorChannel)
		if floor == INVALID_FLOOR && elevatorIsApproaching == false {
			motorChannel <- DIR_DOWN
			LastFloor = INVALID_FLOOR
			elevatorIsApproaching = true
		}
		if floor != INVALID_FLOOR {
			motorChannel <- DIR_STOP
			LastFloor = floor
			state = IDLE
			break
		}
	}
	isInitialized = true
	previousState = INIT
}

func idle(motorChannel chan int, lightChannel chan Light, doorChannel chan bool) {
	if previousState != IDLE {
		fmt.Println("Idle!!!")
		if TargetFloor == UNDEFINED_TARGET_FLOOR {
			// Send network message to recieve TargetFloor
			isOrderServed = false
		} else if TargetFloor > LastFloor {
			state = UP
		} else if TargetFloor < LastFloor {
			state = DOWN
		} else {
			if !isOrderServed { // -- That is TargetFloor == LastFloor
				doorChannel <- OpenDoor
				isOrderServed = true
			}
			state = IDLE
		}
	}
	previousState = IDLE
}

func up(motorChannel chan int, floorChannel chan int) {
	fmt.Println("up!!!")
	floor := ReadFloorSensor(floorChannel)
	if floor == MAX_FLOOR_NUMBER {
		state = DOWN
		return
	}
	motorChannel <- DIR_UP
	state = FLOOR_UP
	previousState = UP
}

func down(motorChannel chan int, floorChannel chan int) {
	fmt.Println("down!!!")
	floor := ReadFloorSensor(floorChannel)
	if floor == 0 {
		state = UP
		return
	}
	motorChannel <- DIR_DOWN
	state = FLOOR_DOWN
	previousState = DOWN
}

func floorUp(motorChannel chan int, lightChannel chan Light, floorChannel chan int) {
	fmt.Println("floor_up!!!")
	floor := ReadFloorSensor(floorChannel)
	if (LastFloor < 0) || (floor > MAX_FLOOR_NUMBER) {
		fmt.Println("Last floor: ", LastFloor)
		fmt.Println("Floor Error")
	}
	switch floor {
	case INVALID_FLOOR:
		fmt.Println("At invalid")
		state = FLOOR_UP
	case TargetFloor:
		fmt.Println("At target")
		motorChannel <- DIR_STOP
		UpdateFloorIndicator(floor, LastFloor, lightChannel)
		UpdateIndicator(BUTTON_HALL_UP, false, floor, lightChannel)
		time.Sleep(1 * time.Second)
		LastFloor = floor
		state = IDLE
	default: // Valid floor which isn't target floor!
		fmt.Println("At default")
		UpdateFloorIndicator(floor, LastFloor, lightChannel)
		motorChannel <- DIR_STOP
		time.Sleep(1 * time.Second)
		motorChannel <- DIR_UP
		if floor > MAX_FLOOR_NUMBER {
			state = FLOOR_DOWN
		} else {
			state = FLOOR_UP
		}
	}
	previousState = FLOOR_UP
}

func floorDown(motorChannel chan int, lightChannel chan Light, floorChannel chan int) {
	fmt.Println("floor_down!!!")
	floor := ReadFloorSensor(floorChannel)
	if (LastFloor < 0) || (floor > MAX_FLOOR_NUMBER) {
		fmt.Println("Last floor: ", LastFloor)
		fmt.Println("Floor Error")
	}
	switch floor {
	case INVALID_FLOOR:
		fmt.Println("At invalid")
		state = FLOOR_DOWN
	case TargetFloor:
		fmt.Println("At target")
		motorChannel <- DIR_STOP
		UpdateFloorIndicator(floor, LastFloor, lightChannel)
		UpdateIndicator(BUTTON_HALL_DOWN, false, floor, lightChannel)
		time.Sleep(1 * time.Second)
		LastFloor = floor
		state = IDLE
	default: // Valid floor which isn't target floor!
		fmt.Println("At default")
		motorChannel <- DIR_STOP
		time.Sleep(1 * time.Second)
		UpdateFloorIndicator(floor, LastFloor, lightChannel)
		if floor <= 0 {
			state = FLOOR_UP
		} else {
			state = FLOOR_DOWN
		}
	}
	previousState = FLOOR_DOWN
}

func FiniteStateMachine(motorChannel chan int, lightChannel chan Light, floorChannel chan int, doorChannel chan bool) {
	for {
		if !isInitialized {
			state = INIT
		}
		switch state {
		case INIT:
			initState(motorChannel, lightChannel, floorChannel)
		case IDLE:
			idle(motorChannel, lightChannel, doorChannel)
		case UP:
			up(motorChannel, floorChannel)
		case FLOOR_UP:
			floorUp(motorChannel, lightChannel, floorChannel)
		case DOWN:
			down(motorChannel, floorChannel)
		case FLOOR_DOWN:
			floorDown(motorChannel, lightChannel, floorChannel)
		}
	}
}
