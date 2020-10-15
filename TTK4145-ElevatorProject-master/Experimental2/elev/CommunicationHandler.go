package elev

import (
	"fmt"
	"strconv"
	"time"
)

type HallOrderElement struct {
	Command   MessageEvent
	Direction MotorDirection // Was string before!
	Floor     int
	Status    int
	ReserveID string
}

var HallOrderTable []HallOrderElement

// In use
func NewOrderEvent(message ElevatorOrderMessage, sendChannel chan ElevatorOrderMessage) { // Master recv
	if message.Origin != "node" {
		return
	}
	tableElement := createTableElement(message)
	if !isElementInHallTable(tableElement) {
		HallOrderTable = append(HallOrderTable, tableElement)
	}
	sendChannel <- ElevatorOrderMessage{
		Event:     EVENT_ACK_NEW_ORDER,
		Direction: message.Direction,
		Floor:     message.Floor,
		Origin:    message.Origin,
		Sender:    "master",
	}
}

// In use
func AckNewOrderEvent(message ElevatorOrderMessage, lightChannel chan Light) { // Node recv
	var lightButton int
	switch message.Direction {
	case DIR_UP:
		lightButton = BUTTON_HALL_UP
	case DIR_DOWN:
		lightButton = BUTTON_HALL_DOWN
	}
	lightChannel <- Light{
		LightType:   lightButton,
		LightOn:     true,
		FloorNumber: message.Floor,
	}
}

// In Use
func OrderReserveEvent(message ElevatorOrderMessage, sendChannel chan ElevatorOrderMessage) { // Master recv
	if message.Origin != "node" {
		return
	}
	/*
		var nextFloor int
		nextFloor = GetOrder(message.Floor, message.Direction)
		SetOrderStatus(STATUS_OCCUPIED, message.Origin, nextFloor, message.Direction)
		sendChannel <- ElevatorOrderMessage{
			Event:      EVENT_ACK_ORDER_RESERVE,
			Floor:      nextFloor,
			AssignedTo: "origin",
			Origin:     "master",
			Sender:     "master",
		}
	*/
	nextFloor := UNDEFINED
	if HallOrderAbove(message.Floor) && HallOrderBelow(message.Floor) {
		floorAbove := GetHallOrderAbove(LastFloor)
		floorBelow := GetHallOrderBelow(LastFloor)
		distanceAbove := LastFloor - floorAbove
		distanceBelow := floorBelow - LastFloor
		fmt.Println("Dist above: ", distanceAbove)
		fmt.Println("Dist under: ", distanceBelow)
		if distanceBelow <= distanceAbove {
			nextFloor = floorBelow
		} else {
			nextFloor = floorAbove
		}
	} else if HallOrderAbove(message.Floor) {
		nextFloor = GetHallOrderAbove(message.Floor)
	} else if HallOrderBelow(message.Floor) {
		nextFloor = GetHallOrderBelow(message.Floor)
	} else {
		nextFloor = message.Floor
	}
	sendChannel <- ElevatorOrderMessage{
		Event:      EVENT_ACK_ORDER_RESERVE,
		Floor:      nextFloor,
		AssignedTo: "origin",
		Origin:     "master",
		Sender:     "master",
	}
}

// In use
func AckOrderReserveEvent(message ElevatorOrderMessage) { // Node recv
	if message.Origin == "master" {
		if message.Floor != UNDEFINED {
			//fmt.Println("TargetFloor ack: " + strconv.Itoa(message.Floor))
			TargetFloor = message.Floor
		}
		//if message.Floor != UNDEFINED {
		//ReserveTable = append(ReserveTable, ReserveElement{Floor: message.Floor})
		//}

	}
}

// In use
func OrderReserveSpecificEvent(message ElevatorOrderMessage, sendChannel chan ElevatorOrderMessage) {
	if message.Origin != "node" {
		return
	}
	if IsOrderAt(message.Floor, message.Direction) {
		sendChannel <- ElevatorOrderMessage{
			Event:      EVENT_ACK_ORDER_RESERVE_SPECIFIC,
			Floor:      message.Floor,
			AssignedTo: "origin",
			Origin:     "master",
			Sender:     "master",
		}
	} else {
		sendChannel <- ElevatorOrderMessage{
			Event:      EVENT_ACK_ORDER_RESERVE_SPECIFIC,
			Floor:      UNDEFINED,
			AssignedTo: "origin",
			Origin:     "master",
			Sender:     "master",
		}
	}
}

// In use
func AckOrderReserveSpecificEvent(message ElevatorOrderMessage) {
	if message.Origin == "master" {
		//if message.Floor != UNDEFINED {
		//AddReservationOrder(message.Floor)
		//}
		if message.Floor != UNDEFINED {
			IsIntermediateStop = true
		} else {
			IsIntermediateStop = false
		}
	}
}

// In Use
func OrderDoneEvent(message ElevatorOrderMessage, sendChannel chan ElevatorOrderMessage) {
	fmt.Println("Order Done")
	if message.Origin != "node" {
		return
	}
	RemoveHallOrder(message.Floor)
	sendChannel <- ElevatorOrderMessage{
		Event:      EVENT_ACK_ORDER_DONE,
		Floor:      message.Floor,
		AssignedTo: "origin",
		Origin:     "master",
		Sender:     "master",
	}
}

// In Use
func AckOrderDoneEvent(message ElevatorOrderMessage, lightChannel chan Light) {
	if message.Origin == "master" {
		lightChannel <- Light{
			LightType:   BUTTON_HALL_UP,
			LightOn:     false,
			FloorNumber: message.Floor,
		}
		lightChannel <- Light{
			LightType:   BUTTON_HALL_DOWN,
			LightOn:     false,
			FloorNumber: message.Floor,
		}
	}
}

// In Use
func createTableElement(message ElevatorOrderMessage) HallOrderElement {
	tableElement := HallOrderElement{
		Command:   message.Event,
		Direction: message.Direction,
		Floor:     message.Floor,
		ReserveID: message.Origin,
		Status:    STATUS_AVAILABLE,
	}
	return tableElement
}

// In use
func isElementInHallTable(element HallOrderElement) bool {
	for _, tableElement := range HallOrderTable {
		if isTableElementEqual(element, tableElement) {
			return true
		}
	}
	return false
}

// In use
func isTableElementEqual(element HallOrderElement, tableElement HallOrderElement) bool {
	if element.Command == tableElement.Command && element.Direction == tableElement.Direction && element.Floor == tableElement.Floor {
		return true
	}
	return false
}

func UpdateReservationTable(sendChannel chan ElevatorOrderMessage) {
	if len(CabOrderTable) != 0 {
		for _, cabOrder := range CabOrderTable {
			ReserveTable = append(ReserveTable, ReserveElement{Floor: cabOrder.Floor})
			removeCabOrder(cabOrder)
		}
	}
	sendChannel <- ElevatorOrderMessage{
		Event:     EVENT_ORDER_RESERVE_SPECIFIC,
		Direction: ElevatorDirection,
		Floor:     LastFloor,
		Origin:    "node",
		Sender:    "node",
	}
}

func CheckForOrders(sendChannel chan ElevatorOrderMessage) {
	if len(CabOrderTable) != 0 {
		for _, cabOrder := range CabOrderTable {
			//TargetFloor = cabOrder.Floor // Fix ordering Problem! Sort based on direction!
			TargetFloor = GetCabOrder(LastFloor, ElevatorDirection)
			removeCabOrder(cabOrder)
			break
		}
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

func removeCabOrder(cabOrder CabOrderElement) {
	for index, element := range CabOrderTable {
		if element == cabOrder {
			CabOrderTable = append(CabOrderTable[:index], CabOrderTable[index+1:]...)
		}
	}
}

func PrintHallTable() {
	if len(HallOrderTable) == 0 {
		fmt.Println("No hall")
	} else {

		fmt.Println("-------------------------------")
		for _, tableElement := range HallOrderTable {
			switch tableElement.Direction {
			case DIR_UP:
				if tableElement.Status == STATUS_AVAILABLE {
					fmt.Println("Hall-Table: " + string(tableElement.Command) + " " + "UP" + " " + strconv.Itoa(tableElement.Floor) + " " + tableElement.ReserveID + " " + "AVAILABLE")
				}
				if tableElement.Status == STATUS_OCCUPIED {
					fmt.Println("Hall-Table: " + string(tableElement.Command) + " " + "UP" + " " + strconv.Itoa(tableElement.Floor) + " " + tableElement.ReserveID + " " + "OCCUPIED")
				}

			case DIR_DOWN:
				if tableElement.Status == STATUS_AVAILABLE {
					fmt.Println("Hall-Table: " + string(tableElement.Command) + " " + "DOWN" + " " + strconv.Itoa(tableElement.Floor) + " " + tableElement.ReserveID + " " + "AVAILABLE")
				}
				if tableElement.Status == STATUS_OCCUPIED {
					fmt.Println("Hall-Table: " + string(tableElement.Command) + " " + "DOWN" + " " + strconv.Itoa(tableElement.Floor) + " " + tableElement.ReserveID + " " + "OCCUPIED")
				}
				fmt.Println("-------------------------------")
			}
		}
	}
}

func PrintCabTable() {
	if len(CabOrderTable) == 0 {
		fmt.Println("No cab")
	} else {
		fmt.Println("-----------------Cab Order------------------")
		for _, tableElement := range CabOrderTable {
			fmt.Println("Order at Floor: " + strconv.Itoa(tableElement.Floor))
			fmt.Println("-------------------------------------------")
		}
	}
}

func PrintReservedTable() {
	if len(ReserveTable) == 0 {
		fmt.Println("No Reservations")
	} else {
		fmt.Println("-----------------Reserved Order------------------")
		for _, tableElement := range ReserveTable {
			fmt.Println("Target at Floor: " + strconv.Itoa(tableElement.Floor))
			fmt.Println("-------------------------------------------")
		}
	}
}

func PrintElevatorInfo() {
	for {
		fmt.Println("------------Elevator Info --------------------")
		fmt.Println("Target floor: " + strconv.Itoa(TargetFloor))
		fmt.Println("Last floor: " + strconv.Itoa(LastFloor))
		fmt.Println("Is intermediate stop: " + strconv.FormatBool(IsIntermediateStop))
		switch ElevatorDirection {
		case DIR_STOP:
			fmt.Println("Direction: DIR_STOP")
		case DIR_UP:
			fmt.Println("Direction: DIR_UP")
		case DIR_DOWN:
			fmt.Println("Direction: DIR_DOWN")
		}
		PrintCabTable()
		PrintHallTable()
		PrintReservedTable()
		time.Sleep(2 * time.Second)
	}
}
