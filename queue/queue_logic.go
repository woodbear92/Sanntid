package queue

import (
	"../config"
	"../elevator_control/elevio"
	"../elevator_control/logic"
	"../typedef"
	"math"
)

const _numElevators int = config.NumElevators
const _numFloors int = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons
const highestPossibleCost int = math.MaxInt64

func TimeToIdle(elevator typedef.ElevatorData) int {

	duration := 0

	if elevator.Connected == false || elevator.HardwareFailure == true {
		return highestPossibleCost
	}

	if logic.OrdersEmpty(elevator) {
		return duration
	}

	switch elevator.State {
	case typedef.State_Idle:
		elevator.CurrentDirection = logic.FindDirection(elevator)
		if elevator.CurrentDirection == elevio.MD_Stop {
			return duration
		}
	case typedef.State_Moving:
		duration += int(config.TravelTimeMs / 2)
		elevator.Floor += int(elevator.CurrentDirection)
	case typedef.State_DoorOpen:
		duration += int(config.DoorOpenTimeMs / 2)
	}

	for {

		if logic.ShouldStop(elevator) {
			elevator.LocalOrders = logic.ExecuteOrder(elevator)
			duration += int(config.DoorOpenTimeMs)
		}

		if logic.OrdersEmpty(elevator) {
			return duration
		}

		elevator.CurrentDirection = logic.FindDirection(elevator)
		if elevator.CurrentDirection != elevio.MD_Stop {
			elevator.State = typedef.State_Moving
		}

		elevator.Floor += int(elevator.CurrentDirection)
		duration += int(config.TravelTimeMs)
	}
}

/*
Below is a implementation of TimeToServeRequest cost algorithm, still need some work
*/
/*
const TRAVEL_TIME=int(config.ElevatorTravelTimeMs)
const DOOR_OPEN_TIME=int(config.DoorOpenTimeMs)
const _numElevators = config.NumElevators
const _numFloors = config.NumFloors

func TimeToServeRequest(elevator typedef.ElevatorData, button elevio.ButtonType, floor int) int{
  var e typedef.ElevatorData
  e=elevator
  e.LocalOrders[button][floor]=typedef.Order_Handle
  duration:=0
  switch e.State {
    case typedef.State_Idle:
      e.CurrentDirection=logic.FindDirection(e)
      if e.CurrentDirection==elevio.MD_Stop {
        return duration
      }
    case typedef.State_Moving:
      duration+=TRAVEL_TIME/2
      e.CurrentDirection=logic.FindDirection(e)
      e.Floor+=int(e.CurrentDirection)
    case typedef.State_DoorOpen:
      duration-=DOOR_OPEN_TIME/2
    default:
      fmt.Println("State undefined")
      duration=math.MaxInt64
      return duration

  }
  for{
    //fmt.Println(e.CurrentDirection)
    if(logic.ShouldStop(e)){
      e.LocalOrders=logic.ExecuteOrder(e)
      if(e.Floor==floor){
        return duration
      }
      duration+=DOOR_OPEN_TIME
      e.CurrentDirection=logic.FindDirection(e)

    }
    e.Floor+=int(e.CurrentDirection)

    duration+=TRAVEL_TIME
    }
}

func ClearOrdersAtCurrentFloor(e typedef.ElevatorData)typedef.ElevatorData{
  e2:=copyElevator(e)
  for button:=0;button<3; button++{
    if e2.LocalOrders[button][e2.Floor]!=typedef.Order_Empty {
      e2.LocalOrders[button][e2.Floor]=typedef.Order_Empty
    }
  }
  return e2;
}
func ChooseElevator(elevators [_numElevators]typedef.ElevatorData, button elevio.ButtonType, floor int) int{
  min:=0
  var elevator_number int
  for i,elevator:=range elevators{
    temp:=TimeToServeRequest(elevator, button, floor)
    if temp==math.MaxInt64{
      continue
    }
    if  i==0 ||temp<min {
      min=temp
      elevator_number=elevator.Number
    }
  }
  return elevator_number
}
func copyElevator(e typedef.ElevatorData)typedef.ElevatorData{
  var e2 typedef.ElevatorData
  e2.HardwareFailure=e.HardwareFailure
  e2.Number=e.Number
  e2.Floor=e.Floor
  e2.Connected=e.Connected
  e2.CurrentDirection=e.CurrentDirection
  e2.LocalOrders=e.LocalOrders
  e2.State=e.State
  return e2
}*/
