package elev

import (
	"fmt"
)

var isInitialized bool = false
var isOrderServed bool = false

var state int
var previousState int

func initState(motorChannel chan MotorDirection, lightChannel chan Light, floorChannel chan int, requestChannel chan Action) {
	fmt.Println("Init!!!")
	if isInitialized {
		state = IDLE
	}
	requestChannel <- Action{
		Command: ACTION_RESET_ALL_LIGHTS,
		Floor:   LastFloor,
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
			ElevatorDirection = DIR_STOP
			LastFloor = floor
			state = IDLE
			break
		}
	}
	isInitialized = true
	previousState = INIT
}

func idle(motorChannel chan MotorDirection, lightChannel chan Light, doorChannel chan bool, requestChannel chan Action, floorChannel chan int) {
	//if previousState != IDLE {
	//ElevatorDirection = DIR_STOP
	if previousState != IDLE { // On new floor enter
		fmt.Println("Idle!!!")
	}
	InitOrderHandeler(lightChannel, doorChannel, requestChannel)
	SetTargetFloor(doorChannel)
	//}
	previousState = IDLE
}

func up(motorChannel chan MotorDirection, floorChannel chan int) {
	//fmt.Println("up!!!")
	floor := ReadFloorSensor(floorChannel)
	if floor == MAX_FLOOR_NUMBER {
		state = DOWN
		return
	}
	ElevatorDirection = DIR_UP
	motorChannel <- DIR_UP
	state = FLOOR_UP
	previousState = UP
}

func down(motorChannel chan MotorDirection, floorChannel chan int) {
	//fmt.Println("down!!!")
	floor := ReadFloorSensor(floorChannel)
	if floor == 0 {
		state = UP
		return
	}
	ElevatorDirection = DIR_DOWN
	motorChannel <- DIR_DOWN
	state = FLOOR_DOWN
	previousState = DOWN
}

func floorUp(motorChannel chan MotorDirection, lightChannel chan Light, floorChannel chan int, doorChannel chan bool, requestChannel chan Action) {
	//fmt.Println("floor_up!!!")
	floor := ReadFloorSensor(floorChannel)
	if (LastFloor < 0) || (floor > MAX_FLOOR_NUMBER) {
		fmt.Println("Last floor: ", LastFloor)
		fmt.Println("Floor Error")
	}
	switch floor {
	case INVALID_FLOOR:
		//fmt.Println("At invalid")
		state = FLOOR_UP
	case TargetFloor:
		TargetFloorUp(floor, LastFloor, lightChannel, motorChannel, requestChannel)
	default: // Valid floor which isn't target floor!
		NonTargetFloorUp(floor, LastFloor, doorChannel, lightChannel, motorChannel, requestChannel)
	}
	previousState = FLOOR_UP
}

func floorDown(motorChannel chan MotorDirection, lightChannel chan Light, floorChannel chan int, doorChannel chan bool, requestChannel chan Action) {
	//fmt.Println("floor_down!!!")
	floor := ReadFloorSensor(floorChannel)
	if (LastFloor < 0) || (floor > MAX_FLOOR_NUMBER) {
		fmt.Println("Last floor: ", LastFloor)
		fmt.Println("Floor Error")
	}
	switch floor {
	case INVALID_FLOOR:
		//fmt.Println("At invalid")
		state = FLOOR_DOWN
	case TargetFloor:
		TargetFloorDown(floor, LastFloor, doorChannel, lightChannel, motorChannel, requestChannel)
	default: // Valid floor which isn't target floor!
		NonTargetFloorDown(floor, LastFloor, doorChannel, lightChannel, motorChannel, requestChannel)
	}
	previousState = FLOOR_DOWN
}

func FiniteStateMachine(motorChannel chan MotorDirection, lightChannel chan Light, floorChannel chan int, doorChannel chan bool, requestChannel chan Action) {
	for {
		if !isInitialized {
			state = INIT
		}
		switch state {
		case INIT:
			initState(motorChannel, lightChannel, floorChannel, requestChannel)
		case IDLE:
			idle(motorChannel, lightChannel, doorChannel, requestChannel, floorChannel)
		case UP:
			up(motorChannel, floorChannel)
		case FLOOR_UP:
			floorUp(motorChannel, lightChannel, floorChannel, doorChannel, requestChannel)
		case DOWN:
			down(motorChannel, floorChannel)
		case FLOOR_DOWN:
			floorDown(motorChannel, lightChannel, floorChannel, doorChannel, requestChannel)
		}
	}
}
