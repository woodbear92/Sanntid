package logic

import (
	"../../config"
	"../../typedef"
	"../elevio"
)

const _numFloors int = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons

func ShouldStop(e typedef.ElevatorData) bool {
	floor := e.Floor

	switch e.State {
	case typedef.State_Moving:
		if e.CurrentDirection == elevio.MD_Up {
			if e.LocalOrders[elevio.BT_Cab][floor] == typedef.Order_Handle || e.LocalOrders[elevio.BT_HallUp][floor] == typedef.Order_Handle || (e.LocalOrders[elevio.BT_HallDown][floor] == typedef.Order_Handle && !ordersAbove(e, floor)) {
				return true
			}
		}
		if e.CurrentDirection == elevio.MD_Down {
			if e.LocalOrders[elevio.BT_Cab][floor] == typedef.Order_Handle || e.LocalOrders[elevio.BT_HallDown][floor] == typedef.Order_Handle || (e.LocalOrders[elevio.BT_HallUp][floor] == typedef.Order_Handle && !ordersBelow(e, floor)) {
				return true
			}
		}
	case typedef.State_Idle:
		fallthrough
	case typedef.State_DoorOpen:
		if e.LocalOrders[elevio.BT_Cab][floor] == typedef.Order_Handle || e.LocalOrders[elevio.BT_HallDown][floor] == typedef.Order_Handle || e.LocalOrders[elevio.BT_HallUp][floor] == typedef.Order_Handle {
			return true
		}
	default:
		return false
	}
	return false
}

func FindDirection(e typedef.ElevatorData) elevio.MotorDirection {
	if OrdersEmpty(e) {
		return elevio.MD_Stop
	}
	switch e.CurrentDirection {
	case elevio.MD_Stop:
		if ordersAbove(e, e.Floor) {
			return elevio.MD_Up
		} else if ordersBelow(e, e.Floor) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Up:
		if ordersAbove(e, e.Floor) {
			return elevio.MD_Up
		} else if ordersBelow(e, e.Floor) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		if ordersBelow(e, e.Floor) {
			return elevio.MD_Down
		} else if ordersAbove(e, e.Floor) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}

func AddOrder(e typedef.ElevatorData, b elevio.ButtonType, f int) [_numOrderButtons][_numFloors]int {
	switch b {
	case elevio.BT_HallUp:
		e.LocalOrders[b][f] = typedef.Order_Received
	case elevio.BT_HallDown:
		e.LocalOrders[b][f] = typedef.Order_Received
	case elevio.BT_Cab:
		e.LocalOrders[b][f] = typedef.Order_Handle
	}
	return e.LocalOrders
}

func ExecuteOrder(e typedef.ElevatorData) [_numOrderButtons][_numFloors]int {
	for buttonType := 0; buttonType < _numOrderButtons; buttonType++ {
		e.LocalOrders[buttonType][e.Floor] = typedef.Order_Executed
	}
	return e.LocalOrders
}

//Set elevator lamps based on orders
func SetLampsOrders(elevator typedef.ElevatorData) {
	for f := 0; f < _numFloors; f++ {
		for b := 0; b < _numOrderButtons; b++ {
			if elevator.LocalOrders[b][f] >= typedef.Order_LightOn {
				elevio.SetButtonLamp(elevio.ButtonType(b), f, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
			}
		}
	}
}

//Help function to determine next direction
//Check if orders above
func ordersAbove(e typedef.ElevatorData, floor int) bool {
	if floor == _numFloors-1 {
		return false
	}
	for b := 0; b < _numOrderButtons; b++ {
		for f := floor + 1; f < _numFloors; f++ {
			if e.LocalOrders[b][f] == typedef.Order_Handle {
				return true
			}
		}
	}
	return false
}

//check if ordersbelow
func ordersBelow(e typedef.ElevatorData, floor int) bool {
	if floor == 0 {
		return false
	}

	for b := 0; b < 3; b++ {
		for f := 0; f < floor; f++ {
			if e.LocalOrders[b][f] == typedef.Order_Handle {
				return true
			}
		}
	}

	return false
}

//Check if no orders above or below
func OrdersEmpty(e typedef.ElevatorData) bool {
	return !ordersAbove(e, 0) && !ordersBelow(e, _numFloors-1)
}
