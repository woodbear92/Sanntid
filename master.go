package main

import (
	"./network/bcast"
	"fmt"
	"time"
	 "os"
	 "strconv"
	 "./typedef"
	 "./network"
	 "./config"
	 "./master"
)

const _numElevators    int = config.NumElevators
const _numFloors       int = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons
const highestPriority int = _numElevators
func main(){
	elev := os.Args[1:][0]
	elevNum, err := strconv.Atoi(elev)
	if err != nil || elevNum < 1{
		fmt.Println("Input error -- Make sure your elevator number is a positive integer")
		os.Exit(1)
	}
	InitNewMaster(elevNum)
}

func InitNewMaster(elevNum int) {

	AllElevators   := [_numElevators] typedef.ElevatorData{}
	fmt.Println(elevNum)
	masterData := typedef.MasterData{}
	masterData.State = typedef.MasterState_Master
	for i := 0;i<_numElevators;i++{
		masterData.Priorities[i] = highestPriority
	}
	masterData.Priorities[elevNum-1]=-1
	typedef.PrintMasterData(masterData)
	//Channels
	TransmitToElevators     			:= make(chan [_numElevators] typedef.ElevatorData)
	ReceiveFromElevators     			:= make(chan typedef.ElevatorData)
	TransmitBackupElevatorData 							:= make(chan [_numElevators] typedef.ElevatorData)
	TransmitBackupBackupData							:= make(chan typedef.MasterData)
	MasterPingTx 					:= make(chan int)
	elevatorDisconnectedChan  := make(chan int)
	backupDisconnectedChan  := make(chan int)
	backupConnectedChan  := make(chan int)


	//Create goroutines to check for elevator data received and check if backups are connected
	for e := 0; e < _numElevators; e++{
		go network.ReceiveDataFromElevators(ReceiveFromElevators,e,elevatorDisconnectedChan)
		if (e != elevNum-1){
			go network.ReceiveBackupConfirmation(e,backupDisconnectedChan,backupConnectedChan)
			}
	}
	go bcast.Transmitter(config.MasterTXPort, TransmitToElevators)
	go bcast.Transmitter(config.MasterPingPort, MasterPingTx)
	go bcast.Transmitter(config.BackupElevatorDataTxPort, TransmitBackupElevatorData)
	go bcast.Transmitter(config.BackupDataTxPort, TransmitBackupBackupData)
	//Send out a signal every .1s
	go func(){
		for {
			MasterPingTx <- elevNum
			time.Sleep(100*time.Millisecond)
		}
	}()
	//Send out to backups every .1s
	go func(){
		for{
			TransmitBackupBackupData <-masterData
			TransmitBackupElevatorData <- AllElevators
			time.Sleep(100*time.Millisecond)
		}
	}()
	//Print states every 2 seconds
	go func(){
		for{
			time.Sleep(2000*time.Millisecond)
				fmt.Println("I am Master")
				typedef.AllElevatorsPrint(AllElevators)
				typedef.PrintMasterData(masterData)
			}

		}()

	for{
		select{
		case e := <-ReceiveFromElevators:
			AllElevators[e.Number] = e
			AllElevators = master.DistributeOrders(AllElevators, e)
			TransmitToElevators <- AllElevators
		case e := <- elevatorDisconnectedChan:
			AllElevators[e].Connected = false
			AllElevators = master.RedistributeOrders(AllElevators, AllElevators[e])
			TransmitToElevators <- AllElevators
		case b := <-backupDisconnectedChan:
			masterData.Priorities[b] = highestPriority
			masterData.Priorities = master.RecalculatePriorities(masterData.Priorities)
		case b := <-backupConnectedChan:
			if (masterData.Priorities[b] == highestPriority){
			masterData.Priorities[b] = highestPriority-1
			masterData.Priorities = master.RecalculatePriorities(masterData.Priorities)
			}

	}
}
}
