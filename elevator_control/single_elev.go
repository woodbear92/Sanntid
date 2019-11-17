package elevator_control

import (
	"./fsm"
	"fmt"
)

func InitializeSingleElevator(elevatorNumber int, portNum string) {
	fmt.Println("elev:", elevatorNumber, "on port", portNum)
	fsm.Fsm(elevatorNumber, portNum)
}
