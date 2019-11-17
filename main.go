package main

import (
	"./config"
	"./elevator_control"
	"./network/bcast"
	"./queue"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	elevNum := os.Args[1:][0]
	portNum := os.Args[1:][1]

	elevatorNumber, err := strconv.Atoi(elevNum)
	if err != nil || elevatorNumber < 1 {
		fmt.Println("Input error -- Make sure your elevator number is a positive integer")
		os.Exit(1)
	}
	MasterPingRx := make(chan int)
	waitForMasterTimer := time.NewTimer(config.BackupWaitTimeMs * time.Millisecond * 2)

	go bcast.Receiver(config.MasterPingPort, MasterPingRx)
l:
	for {
		select {
		case m := <-MasterPingRx:
			//Master exists, spawn as backup
			fmt.Println("Master sent a ping, spawing as backup")
			succeed := queue.Spawn("backup", elevatorNumber, m)
			if !succeed {
				fmt.Println("Spawning failed")
			}
			break l
		case <-waitForMasterTimer.C:
			//Spawn as master
			fmt.Println("Spawning as new master")
			succeed := queue.Spawn("master", elevatorNumber, 0)
			if !succeed {
				fmt.Println("Spawning failed")
			}
			break l
		default:
			break
		}
	}

	elevator_control.InitializeSingleElevator(elevatorNumber, portNum)
}
