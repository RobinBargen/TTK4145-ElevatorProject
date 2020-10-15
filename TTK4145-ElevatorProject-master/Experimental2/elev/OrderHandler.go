package elev

import (
	"fmt"
)

type CabOrderElement struct {
	Floor int
}
type ReserveElement struct {
	Floor int
}

var CabOrderTable []CabOrderElement
var ReserveTable []ReserveElement

// In use
func IsOrderAt(floor int, direction MotorDirection) bool {
	for _, element := range HallOrderTable {
		if element.Floor == floor && element.Direction == direction {
			return true
		}
	}
	return false
}

func SetOrderStatus(status int, id string, floor int, direction MotorDirection) {
	element := HallOrderElement{
		Direction: direction,
		Floor:     floor,
	}
	for _, tableElement := range HallOrderTable {
		if element.Direction == tableElement.Direction && element.Floor == tableElement.Floor {
			tableElement.Status = status
			tableElement.ReserveID = id
		}
	}
}

func AddCabOrder(floor int) {
	if !isElementInCabTable(floor) {
		cabOrder := CabOrderElement{
			Floor: floor,
		}
		CabOrderTable = append(CabOrderTable, cabOrder)
	}
}

func AddReservationOrder(floor int) {
	if !IsElementInReserveTable(floor) {
		order := ReserveElement{
			Floor: floor,
		}
		ReserveTable = append(ReserveTable, order)
	}
}
func IsElementInReserveTable(floor int) bool {
	for _, element := range ReserveTable {
		if element.Floor == floor {
			return true
		}
	}
	return false
}

func checkCabOrderAtFloor(floor int) int {
	for _, cabOrder := range CabOrderTable {
		if floor == cabOrder.Floor {
			return floor
		}
	}
	return UNDEFINED
}

func CabOrderAtFloor(floor int) int {
	for _, cabOrder := range CabOrderTable {
		if floor == cabOrder.Floor {
			return floor
		}
	}
	return UNDEFINED
}

func checkCabOrderAbove(floor int, direction MotorDirection) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	for _, order := range CabOrderTable {
		if direction == DIR_UP {
			distance = order.Floor - floor
			if minDistance == UNDEFINED {
				minDistance = distance
			}
			if distance < minDistance {
				distance = minDistance
				bestFloor = order.Floor
			}
		}
	}
	return bestFloor
}

func checkCabOrderBelow(floor int, direction MotorDirection) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	for _, order := range CabOrderTable {
		if direction == DIR_DOWN {
			distance = floor - order.Floor
			if minDistance == UNDEFINED {
				minDistance = distance
			}
			if distance < minDistance {
				distance = minDistance
				bestFloor = order.Floor
			}
		}
	}
	fmt.Println(bestFloor)
	return bestFloor
}

func GetCabOrder(floor int, direction MotorDirection) int {
	var nextFloor int
	nextFloor = UNDEFINED
	switch direction {
	case DIR_STOP:
		nextFloor = checkCabOrderAtFloor(floor)
	case DIR_UP:
		nextFloor = checkCabOrderAbove(floor, direction)
	case DIR_DOWN:
		nextFloor = checkCabOrderBelow(floor, direction)
	}
	return nextFloor
}

func InternalOrderCheck(sendChannel chan ElevatorOrderMessage, floorChannel chan int) {
	for {
		if TargetFloor != UNDEFINED_TARGET_FLOOR {
			return
		}
		floor := ReadFloorSensor(floorChannel)
		if len(CabOrderTable) != 0 && floor != INVALID_FLOOR {
			TargetFloor = GetCabOrder(floor, ElevatorDirection)
		} else {
			sendChannel <- ElevatorOrderMessage{
				Event:     EVENT_ORDER_RESERVE,
				Direction: ElevatorDirection,
				Floor:     LastFloor,
				Origin:    "node",
				Sender:    "node",
			}
		}
	}
}

func isElementInCabTable(floor int) bool {
	for _, element := range CabOrderTable {
		if element.Floor == floor {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////
// In use
func CabOrderAbove(floor int) bool {
	for _, element := range CabOrderTable {
		if element.Floor > floor {
			return true
		}
	}
	return false
}

// In use
func CabOrderBelow(floor int) bool {
	for _, element := range CabOrderTable {
		if element.Floor < floor {
			return true
		}
	}
	return false
}

// In use
func IsCabFloor(floor int) bool {
	for _, element := range CabOrderTable {
		if element.Floor == floor {
			return true
		}
	}
	return false
}

// In use
func GetCabOrderAbove(floor int) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	minDistance = UNDEFINED
	for _, order := range CabOrderTable {
		distance = order.Floor - floor
		if minDistance == UNDEFINED {
			minDistance = distance
		}
		if bestFloor == UNDEFINED {
			bestFloor = order.Floor
		}
		if distance < minDistance {
			distance = minDistance
			bestFloor = order.Floor
		}
	}
	return bestFloor
}

// In Use
func GetCabOrderBelow(floor int) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	minDistance = UNDEFINED
	for _, order := range CabOrderTable {
		distance = floor - order.Floor
		if minDistance == UNDEFINED {
			minDistance = distance
		}
		if bestFloor == UNDEFINED {
			bestFloor = order.Floor
		}
		if distance < minDistance {
			distance = minDistance
			bestFloor = order.Floor
		}
	}
	return bestFloor
}

// In use
func RemoveCabOrder(floor int) {
	for index, element := range CabOrderTable {
		if element.Floor == floor {
			CabOrderTable = append(CabOrderTable[:index], CabOrderTable[index+1:]...)
		}
	}
}

//////////////////////////////////////////////////////
// In use
func IsHallFloor(floor int, requestChannel chan Action) bool {
	requestChannel <- Action{
		Command:   ACTION_REQUEST_SPECIFIC_ORDER,
		Direction: ElevatorDirection,
		Floor:     floor,
	}
	return IsIntermediateStop
}

// In Use
func HallOrderAbove(floor int) bool {
	for _, element := range HallOrderTable {
		if element.Floor > floor {
			return true
		}
	}
	return false
}

// In use
func HallOrderBelow(floor int) bool {
	for _, element := range HallOrderTable {
		if element.Floor < floor {
			return true
		}
	}
	return false
}

// In use
func GetHallOrderAbove(floor int) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	minDistance = UNDEFINED
	for _, order := range HallOrderTable {
		distance = order.Floor - floor
		if minDistance == UNDEFINED {
			minDistance = distance
		}
		if bestFloor == UNDEFINED {
			bestFloor = order.Floor
		}
		if distance < minDistance {
			distance = minDistance
			bestFloor = order.Floor
		}
	}
	return bestFloor
}

// In use
func GetHallOrderBelow(floor int) int {
	var bestFloor int
	var minDistance int
	var distance int

	bestFloor = UNDEFINED
	minDistance = UNDEFINED
	for _, order := range HallOrderTable {
		distance = floor - order.Floor
		if minDistance == UNDEFINED {
			minDistance = distance
		}
		if bestFloor == UNDEFINED {
			bestFloor = order.Floor
		}
		if distance < minDistance {
			distance = minDistance
			bestFloor = order.Floor
		}
	}
	return bestFloor
}

func RemoveHallOrder(floor int) {
	fmt.Println("Order remove func")
	fmt.Println("Floor: ", floor)
	for index, order := range HallOrderTable {
		if len(HallOrderTable) == 0 || len(HallOrderTable) <= index {
			break
		}
		if order.Floor == floor {
			fmt.Println("Slice size, index: ", len(HallOrderTable), index)
			HallOrderTable = append(HallOrderTable[:index], HallOrderTable[index+1:]...)
		}
	}
}
