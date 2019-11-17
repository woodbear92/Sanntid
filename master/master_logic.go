package master

import (
	"../config"
	"../elevator_control/elevio"
	"../queue"
	"../typedef"
	"math"
)

const _numElevators int = config.NumElevators
const _numFloors int = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons
const highestPriority int = _numElevators

func DistributeOrders(elevators [_numElevators]typedef.ElevatorData, orderFromElevator typedef.ElevatorData) [_numElevators]typedef.ElevatorData {

	elev := orderFromElevator.Number
	toElev := findElevator(elevators, orderFromElevator)

	for b := 0; b < _numOrderButtons; b++ {
		for f := 0; f < _numFloors; f++ {
			for e := 0; e < _numElevators; e++ {
				switch orderFromElevator.LocalOrders[b][f] {
				case typedef.Order_Received:
					for l := 0; l < _numElevators; l++ {
						if l == toElev {
							elevators[l].LocalOrders[b][f] = typedef.Order_Handle
						} else {
							elevators[l].LocalOrders[b][f] = typedef.Order_LightOn
						}
					}
				case typedef.Order_Executed:
					if b == elevio.BT_Cab {
						elevators[elev].LocalOrders[b][f] = typedef.Order_Empty
					} else {
						for l := 0; l < _numElevators; l++ {
							elevators[l].LocalOrders[b][f] = typedef.Order_Empty
						}
					}
				case typedef.Order_Handle:
					if b == elevio.BT_Cab {
						elevators[elev].LocalOrders[b][f] = typedef.Order_Handle
					}
				default:
					break
				}
			}
		}
	}
	return elevators
}

func RedistributeOrders(elevators [_numElevators]typedef.ElevatorData, orderFromElevator typedef.ElevatorData) [_numElevators]typedef.ElevatorData {

	redistributeOrders := orderFromElevator.LocalOrders
	newElev := findElevator(elevators, orderFromElevator)

	for b := 0; b < _numOrderButtons-1; b++ {
		for f := 0; f < _numFloors; f++ {
			if redistributeOrders[b][f] == typedef.Order_Handle {
				for e := 0; e < _numElevators; e++ {
					if e == newElev {
						elevators[e].LocalOrders[b][f] = typedef.Order_Handle
					} else {
						elevators[e].LocalOrders[b][f] = typedef.Order_LightOn
					}
				}
			}
		}
	}
	return elevators
}

func RecalculatePriorities(priorities [_numElevators]int) [_numElevators]int {
	list := priorities
	priority := 1
	for i, elm := range priorities {
		if elm == highestPriority {
			continue
		} else if elm == -1 {
			continue
		} else {
			list[i] = priority
			priority += 1
		}
	}
	return list
}

func findElevator(elevators [_numElevators]typedef.ElevatorData, orderFromElevator typedef.ElevatorData) int {
	toElev := orderFromElevator.Number
	lowestCost := math.MaxInt64

	for e := 0; e < _numElevators; e++ {
		for b := 0; b < _numOrderButtons; b++ {
			for f := 0; f < _numFloors; f++ {
				if orderFromElevator.LocalOrders[b][f] == typedef.Order_Received {
					elevators[e].LocalOrders[b][f] = typedef.Order_Handle
				}
			}
		}
		cost := queue.TimeToIdle(elevators[e])
		if cost < lowestCost {
			toElev = e
			lowestCost = cost
		}
	}
	return toElev
}
