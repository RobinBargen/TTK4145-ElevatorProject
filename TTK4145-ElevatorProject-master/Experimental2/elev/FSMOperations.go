package elev

import (
  "fmt"
  "time"
)
/****************************************************************
* FSM Help functions
****************************************************************/

/*---------------------------------- Idle ----------------------*/
func SetTargetFloor(doorChannel chan bool) {
  if TargetFloor == LastFloor || TargetFloor == UNDEFINED_TARGET_FLOOR {
		if !isOrderServed { // -- That is TargetFloor == LastFloor
			doorChannel <- OpenDoor
			isOrderServed = true
		}
		state = IDLE
	} else if TargetFloor > LastFloor {
		isOrderServed = false
		state = UP
	} else if TargetFloor < LastFloor {
		isOrderServed = false
		state = DOWN
	} else {
		// Do nothing!!
	}
}

func InitOrderHandeler(lightChannel chan Light, doorChannel chan bool, requestChannel chan Action) {
  if IsCabFloor(LastFloor) {
		fmt.Println("Hei 0")
		lightChannel <- Light{
			LightType:   BUTTON_CAB,
			LightOn:     false,
			FloorNumber: LastFloor,
		}
    RemoveCabOrder(LastFloor)
    doorChannel <- true
		time.Sleep(2 * time.Second)
	} else if CabOrderAbove(LastFloor) && CabOrderBelow(LastFloor) {
		fmt.Println("Hei 1")
		floorAbove := GetCabOrderAbove(LastFloor)
		floorBelow := GetCabOrderBelow(LastFloor)
		distanceAbove := LastFloor - floorAbove
		distanceBelow := floorBelow - LastFloor

		if distanceBelow <= distanceAbove {
			TargetFloor = floorBelow
			RemoveCabOrder(floorAbove)
			ServeOrder(doorChannel, requestChannel)
		} else {
			TargetFloor = floorAbove
			RemoveCabOrder(floorBelow)
			ServeOrder(doorChannel, requestChannel)
		}

	} else if CabOrderAbove(LastFloor) {
		fmt.Println("Hei 2")
		newTarget := GetCabOrderAbove(LastFloor)
		TargetFloor = newTarget
		RemoveCabOrder(newTarget)
		ServeOrder(doorChannel, requestChannel)
	} else if CabOrderBelow(LastFloor) {
		fmt.Println("Hei 3")
		newTarget := GetCabOrderBelow(LastFloor)
		TargetFloor = newTarget
		RemoveCabOrder(newTarget)
		ServeOrder(doorChannel, requestChannel)
	} else {
		fmt.Println("Hei 4")
		requestChannel <- Action{
			Command:   ACTION_REQUEST_ORDER,
			Direction: ElevatorDirection,
			Floor:     LastFloor,
		}
		requestChannel <- Action{
			Command: ACTION_ORDER_DONE,
			Floor:   LastFloor,
		}
		time.Sleep(2 * time.Second)
	}
}

/*-------------------------- Floor Up ------------------------*/
func TargetFloorUp(floor int, lastFloor int, lightChannel chan Light, motorChannel chan MotorDirection, requestChannel chan Action) {
  fmt.Println("At target")
  motorChannel <- DIR_STOP
  UpdateFloorIndicator(floor, lastFloor, lightChannel)
  //UpdateIndicator(BUTTON_HALL_UP, false, floor, lightChannel)
  requestChannel <- Action{ // --
    Command: ACTION_ORDER_DONE,
    Floor:   floor,
  } //--
  lightChannel <- Light{
    LightType:   BUTTON_CAB,
    LightOn:     false,
    FloorNumber: floor,
  }
  time.Sleep(500 * time.Millisecond)
  LastFloor = floor
  state = IDLE
}

func NonTargetFloorUp(floor int, lastFloor int, doorChannel chan bool, lightChannel chan Light, motorChannel chan MotorDirection, requestChannel chan Action) {
  fmt.Println("At default")
  LastFloor = floor
  IsIntermediateStop = IsHallFloor(LastFloor, requestChannel)
  UpdateFloorIndicator(floor, LastFloor, lightChannel)
  motorChannel <- DIR_STOP
  time.Sleep(1 * time.Second)
  //////////////////////////
  if IsIntermediateStop {
    doorChannel <- true
    time.Sleep(2 * time.Second)
    IsIntermediateStop = false
    fmt.Println("Try to request order done!!!!")
    requestChannel <- Action{
      Command: ACTION_ORDER_DONE,
      Floor:   LastFloor,
    }
  }
  /////////////////////////
  motorChannel <- DIR_UP
  if floor > MAX_FLOOR_NUMBER {
    state = FLOOR_DOWN
  } else {
    state = FLOOR_UP
  }
}

/*----------------------------------- Floor Down -----------------------------*/
func TargetFloorDown(floor int, lastFloor int, doorChannel chan bool, lightChannel chan Light, motorChannel chan MotorDirection, requestChannel chan Action) {
  fmt.Println("At target")
  motorChannel <- DIR_STOP
  UpdateFloorIndicator(floor, LastFloor, lightChannel)
  //UpdateIndicator(BUTTON_HALL_DOWN, false, floor, lightChannel)
  lightChannel <- Light{
    LightType:   BUTTON_CAB,
    LightOn:     false,
    FloorNumber: floor,
  }
  requestChannel <- Action{ // --
    Command: ACTION_ORDER_DONE,
    Floor:   floor,
  } //--
  time.Sleep(500 * time.Millisecond)
  LastFloor = floor
  state = IDLE
}

func NonTargetFloorDown(floor int, lastFloor int, doorChannel chan bool, lightChannel chan Light, motorChannel chan MotorDirection, requestChannel chan Action) {
  fmt.Println("At default")
  LastFloor = floor
  IsIntermediateStop = IsHallFloor(LastFloor, requestChannel)
  motorChannel <- DIR_STOP
  time.Sleep(1 * time.Second)
  UpdateFloorIndicator(floor, LastFloor, lightChannel)
  //////////////////////////
  if IsIntermediateStop {
    doorChannel <- true
    time.Sleep(2 * time.Second)
    IsIntermediateStop = false
    fmt.Println("Try to request order done!!!!")
    requestChannel <- Action{
      Command: ACTION_ORDER_DONE,
      Floor:   LastFloor,
    }
  }
  /////////////////////////
  motorChannel <- DIR_DOWN
  if floor < 0 {
    state = FLOOR_UP
  } else {
    state = FLOOR_DOWN
  }
}

func ServeOrder(doorChannel chan bool, requestChannel chan Action) {
  requestChannel <- Action{
    Command: ACTION_ORDER_DONE,
    Floor:   LastFloor,
  }
  doorChannel <- true
  time.Sleep(2 * time.Second)
}
