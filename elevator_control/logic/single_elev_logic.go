package logic

import (
	"../../typedef"
	"../elevio"
)

func AddOrderSingleElevator(e typedef.ElevatorData, button elevio.ButtonType, floor int) [_numOrderButtons][_numFloors]int {
	switch button {
	case elevio.BT_Cab:
		e.LocalOrders[button][floor] = typedef.Order_Handle
	case elevio.BT_HallUp:
		fallthrough
	case elevio.BT_HallDown:
		e.LocalOrders[button][floor] = typedef.Order_Handle
	}
	SetLampsOrders(e)
	return e.LocalOrders
}

func ExecuteOrderSingleElevator(e typedef.ElevatorData) [_numOrderButtons][_numFloors]int {
	localOrders := e.LocalOrders
	nextDirection := FindDirection(e)
	for buttonType := 0; buttonType < _numOrderButtons; buttonType++ {
		if localOrders[buttonType][e.Floor] == typedef.Order_Handle {
			if nextDirection == elevio.MD_Up && buttonType == 1 && e.Floor != _numFloors-1 {
				continue
			} else if nextDirection == elevio.MD_Down && buttonType == 0 && e.Floor != 0 {
				continue
			} else {
				localOrders[buttonType][e.Floor] = typedef.Order_Empty
			}
		}
	}
	SetLampsOrders(e)
	return localOrders
}
