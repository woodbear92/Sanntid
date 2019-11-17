package fsm

import (
	"../../config"
	"../../network/bcast"
	"../../typedef"
	"../elevio"
	"../logic"
	"fmt"
	"time"
)

const _numFloors int = config.NumFloors
const _numElevators int = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons

var doorTimer *time.Timer = time.NewTimer(config.DoorOpenTimeMs * time.Millisecond)
var hardwareTimer *time.Timer = time.NewTimer(config.TravelTimeMs * time.Millisecond)

func Fsm(elevatorNumber int, port string) {

	elevio.Init("localhost:"+port, _numFloors)

	//Kill timers
	doorTimer.Stop()
	hardwareTimer.Stop()
	masterPingTimer := time.NewTimer(config.BackupWaitTimeMs * time.Millisecond * 2)

	TransmitToMaster := make(chan typedef.ElevatorData)
	ReceiveFromMaster := make(chan [_numElevators]typedef.ElevatorData)
	masterPingRx := make(chan int)
	buttonPress := make(chan elevio.ButtonEvent)
	floorSensor := make(chan int)

	elevator := typedef.ElevatorData{}

	go bcast.Receiver(config.MasterPingPort, masterPingRx)
	go elevio.PollButtons(buttonPress)
	go elevio.PollFloorSensor(floorSensor)

	elevator = initFSM(elevatorNumber, elevator, floorSensor)

	go bcast.Transmitter(config.ElevatorTXPort+elevator.Number, TransmitToMaster)
	go bcast.Receiver(config.MasterTXPort, ReceiveFromMaster)

	go func() { //Send elevator to master every .1s
		for {
			TransmitToMaster <- elevator
			time.Sleep(100 * time.Millisecond)
			logic.SetLampsOrders(elevator)
		}
	}()

	go func() { //Check that master sends ping
		for {
			select {
			case <-masterPingRx:
				//Master
				masterPingTimer.Reset(config.BackupWaitTimeMs * time.Millisecond * 2)
				elevator.Connected = true
			case <-masterPingTimer.C:
				elevator.Connected = false
			}
		}
	}()

	//The following code is used if obstruction switch og stop button functionality should be required
	/*
		obstructionSwitch     := make(chan bool)
		go elevio.PollObstructionSwitch(obstructionSwitch)

		stopButtonPressed       := make(chan bool)
		go elevio.PollStopButton(stopButtonPressed)
	*/

	for {
		if elevator.HardwareFailure == true {
			fmt.Println("Elevator hardware has failed, trying to restart motor")
		loop:
			for {
				elevio.SetMotorDirection(elevator.CurrentDirection)
				select {
				case f := <-floorSensor:
					fmt.Println("Elevator hardware now functioning")
					elevator.Floor = f
					elevator.HardwareFailure = false
					elevio.SetFloorIndicator(elevator.Floor)
					fmt.Print("Elevator reached floor: ", elevator.Floor, "\n")
					if logic.ShouldStop(elevator) {
						elevator = handleOrder(elevator)
					} else if elevator.Floor == 0 || elevator.Floor == _numFloors-1 {
						elevio.SetMotorDirection(elevio.MD_Stop)
						elevator.State = typedef.State_Idle
						elevator.CurrentDirection = elevio.MD_Stop
					}
					hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond * 2)
					break loop
				case fromMaster := <-ReceiveFromMaster:
					elevatorFromMaster := fromMaster[elevator.Number]
					elevator.LocalOrders = newOrdersFromMaster(elevator, elevatorFromMaster.LocalOrders)
				/*
					case <-obstructionSwitch:
						//Add obstruction switch functionality
					case <-stopButtonPressed
						//Add stop button functionality
				*/
				default:
					time.Sleep(250 * time.Millisecond)
					break
				}
			}
		} else {
			select {
			case b := <-buttonPress:
				if elevator.Connected == false {
					fmt.Println("Elevator not connected, continuing as single elevator")
					elevator.LocalOrders = logic.AddOrderSingleElevator(elevator, b.Button, b.Floor)
				} else {
					elevator.LocalOrders = logic.AddOrder(elevator, b.Button, b.Floor)
				}
				fmt.Println("Elevator button pressed: ", b)
				if logic.ShouldStop(elevator) && elevator.State != typedef.State_Moving { // Order on current floor
					elevator = handleOrder(elevator)
				} else if elevator.State == typedef.State_Idle {
					elevator = setDirection(elevator) // Bestilling når idle
				}
			case f := <-floorSensor:
				elevator.Floor = f
				hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond * 2)
				elevio.SetFloorIndicator(elevator.Floor)
				fmt.Print("Elevator reached floor: ", elevator.Floor, "\n")
				if logic.ShouldStop(elevator) {
					elevator = handleOrder(elevator)
				} else if elevator.Floor == 0 || elevator.Floor == _numFloors-1 {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevator.State = typedef.State_Idle
					elevator.CurrentDirection = elevio.MD_Stop
				}
			case <-doorTimer.C:
				elevio.SetDoorOpenLamp(false)
				elevator = setDirection(elevator)
				if elevator.CurrentDirection == elevio.MD_Stop {
					elevator.State = typedef.State_Idle
				}
				logic.SetLampsOrders(elevator)
			case <-hardwareTimer.C:
				if elevator.State == typedef.State_Moving {
					elevator.HardwareFailure = true
				}
			case fromMaster := <-ReceiveFromMaster:
				elevatorFromMaster := fromMaster[elevator.Number]
				elevator.LocalOrders = newOrdersFromMaster(elevator, elevatorFromMaster.LocalOrders)

				if logic.ShouldStop(elevator) && elevator.State != typedef.State_Moving {
					elevator = handleOrder(elevator)
				} else if elevator.State == typedef.State_Idle {
					elevator = setDirection(elevator)
				}
				/*
					case <-obstructionSwitch:
						//Add obstruction switch functionality
					case <-stopButtonPressed
						//Add stop button functionality
				*/
			}
		}
	}

}

func initFSM(elevatorNumber int, elevator typedef.ElevatorData, floorSensor chan int) typedef.ElevatorData {
	elevio.SetMotorDirection(elevio.MD_Down)
	hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond)
loop:
	for {
		select {
		case <-hardwareTimer.C:
			fmt.Println("Unable to initialize due to hardware failure")
			elevator.HardwareFailure = true
			elevio.SetMotorDirection(elevio.MD_Down)
			hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond)
		case f := <-floorSensor:
			fmt.Println("Floor sensor read with value;", f)
			if f != -1 {
				hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond)
				elevator.Floor = f
				elevator.CurrentDirection = elevio.MD_Stop
				elevator.LocalOrders = [_numOrderButtons][_numFloors]int{}
				elevator.State = typedef.State_Idle
				elevator.Number = elevatorNumber - 1
				elevator.HardwareFailure = false
				elevator.Connected = true
				break loop
			}
		}
	}
	elevio.SetMotorDirection(elevator.CurrentDirection)
	fmt.Println("FSM Initialization complete")
	return elevator
}

func handleOrder(elevator typedef.ElevatorData) typedef.ElevatorData {
	switch elevator.Connected {
	case true:
		elevator.LocalOrders = logic.ExecuteOrder(elevator)
	case false:
		elevator.LocalOrders = logic.ExecuteOrderSingleElevator(elevator)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevator.State = typedef.State_DoorOpen
	elevio.SetDoorOpenLamp(true)
	doorTimer.Reset(config.DoorOpenTimeMs * time.Millisecond)
	return elevator
}

func setDirection(elevator typedef.ElevatorData) typedef.ElevatorData {
	elevator.CurrentDirection = logic.FindDirection(elevator)
	elevio.SetMotorDirection(elevator.CurrentDirection)
	if elevator.CurrentDirection != elevio.MD_Stop {
		elevator.State = typedef.State_Moving
		hardwareTimer.Reset(config.TravelTimeMs * time.Millisecond)
	}
	return elevator
}

func newOrdersFromMaster(elevator typedef.ElevatorData, orderFromMaster [_numOrderButtons][_numFloors]int) [_numOrderButtons][_numFloors]int {
	orders := elevator.LocalOrders
	for f := 0; f < _numFloors; f++ {
		for b := 0; b < _numOrderButtons-1; b++ {
			switch orderFromMaster[b][f] {
			case typedef.Order_Empty:
				if orders[b][f] != typedef.Order_Received {
					orders[b][f] = orderFromMaster[b][f]
				}
			case typedef.Order_Received:
				break
			case typedef.Order_LightOn:
				if elevator.HardwareFailure == true {
					orders[b][f] = orderFromMaster[b][f]
					break
				}
				fallthrough
			case typedef.Order_Handle:
				if orders[b][f] != typedef.Order_Executed {
					orders[b][f] = orderFromMaster[b][f]
				}
			}
		}
		if orderFromMaster[int(elevio.BT_Cab)][f] == typedef.Order_Handle {
			orders[int(elevio.BT_Cab)][f] = orderFromMaster[int(elevio.BT_Cab)][f]
		}
	}
	return orders
}
