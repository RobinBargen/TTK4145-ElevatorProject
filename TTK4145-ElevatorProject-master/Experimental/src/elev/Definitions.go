package elev

type ActionCommand string

type Action struct {
	Command    ActionCommand
	Parameters string
}

type Light struct {
	LightType   int
	LightOn     bool
	FloorNumber int
}

const (
	INIT       int = 0
	IDLE       int = 1
	UP         int = 2
	FLOOR_UP   int = 3
	DOWN       int = 4
	FLOOR_DOWN int = 5
)

const (
	DIR_UP   int = 1
	DIR_DOWN int = -1
	DIR_STOP int = 0
)

const (
	DOOR_CLOSED int = 0
	DOOR_OPEN   int = 1
)

const (
	BUTTON_HALL_UP   int = 0
	BUTTON_HALL_DOWN int = 1
	BUTTON_CAB       int = 2
	FLOOR_INDICATOR  int = 3
)

const (
	COMMAND_ORDER         ActionCommand = "ORDR"
	COMMAND_REQUEST_ORDER ActionCommand = "REQ_ORDR"
	COMMAND_ORDER_DONE    ActionCommand = "ORDR_DONE"
)

const UNDEFINED int = -1
const UNDEFINED_TARGET_FLOOR int = -1
const INVALID_FLOOR int = -1
const MAX_FLOOR_NUMBER int = 4

const (
	CloseDoor bool = false
	OpenDoor  bool = true
)

var LastFloor int
var TargetFloor int

func ReadFloorSensor(floorChannel chan int) int {
	select {
	case floor := <-floorChannel:
		return floor
	default:
		return INVALID_FLOOR
	}
}

func UpdateFloorIndicator(floorNumber int, prevFloorNumber int, lightChannel chan Light) {
	lightChannel <- Light{LightType: FLOOR_INDICATOR, LightOn: false, FloorNumber: prevFloorNumber}
	lightChannel <- Light{LightType: FLOOR_INDICATOR, LightOn: true, FloorNumber: floorNumber}
}
func UpdateIndicator(indicator int, lightActive bool, floorNumber int, lightChannel chan Light) {
	lightChannel <- Light{LightType: indicator, LightOn: lightActive, FloorNumber: floorNumber}
}
