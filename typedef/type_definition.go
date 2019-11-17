package typedef

import (
	"../config"
	"../elevator_control/elevio"
	"fmt"
)

const _numFloors int = config.NumFloors
const _numElevators int = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons

//Structs:
type MasterData struct {
	State      MasterState
	Priorities [_numElevators]int
}

type ElevatorData struct {
	Number           int
	Floor            int
	CurrentDirection elevio.MotorDirection
	LocalOrders      [_numOrderButtons][_numFloors]int
	State            ElevatorState
	HardwareFailure  bool
	Connected        bool
}

type AllElevatorData struct {
	Elevators [_numElevators]ElevatorData
}

//States for single elevator
type ElevatorState int

const (
	State_Idle     = 0
	State_Moving   = 1
	State_DoorOpen = 2
)

//States for master/backups
type MasterState int

const (
	MasterState_Master = 1
	MasterState_Backup = 0
)

type OrderStatus int

const (
	Order_Empty    = 0
	Order_Executed = -1
	Order_Received = 1
	Order_LightOn  = 2
	Order_Handle   = 3
)

//Printing functions below this point.
func PrintMasterData(m MasterData) {
	n := len(m.Priorities)
	fmt.Println("****MASTER DATA****")
	fmt.Println("State:", m.State)
	fmt.Print("Elevator:    ")
	for i := 0; i < n; i++ {
		fmt.Print("|", i)
	}
	fmt.Println("|\nPriorities: ", m.Priorities)
	fmt.Println("*******************")
}

func PrintElevatorData(e ElevatorData) {
	fmt.Println("****ELEVATOR DATA****")
	fmt.Println("Number: ", e.Number)
	fmt.Println("Floor: ", e.Floor)
	fmt.Println("Direction: ", e.CurrentDirection)
	fmt.Println("Orders: ", e.LocalOrders)
	fmt.Println("State: ", stateToString(e.State))
	fmt.Println("Hardware: ", e.HardwareFailure)
	fmt.Println("Connected: ", e.Connected)
	fmt.Println("*********************")
}

func AllElevatorsPrint(Elevators [_numElevators]ElevatorData) {
	fmt.Println("             |Hall up| |Hall dn| |Cab    |")
	for e := 0; e < _numElevators; e++ {
		fmt.Print("Elevator ", e, "  ")
		printElevator(Elevators[e])
	}
}

func printElevator(e ElevatorData) {
	fmt.Print(e.LocalOrders)
	fmt.Print("  N: ", e.Number)
	fmt.Print("  F: ", e.Floor)
	fmt.Print("  D: ", e.CurrentDirection)
	fmt.Print("  S: ", stateToString(e.State))
	fmt.Print("  C: ", e.Connected)
	fmt.Print("  H: ", e.HardwareFailure,"\n")
}

func stateToString(s ElevatorState) string {
	switch s {
	case State_Idle:
		return "Idle"
	case State_Moving:
		return "Moving"
	case State_DoorOpen:
		return "Door open"
	}
	return "Undefined"
}
