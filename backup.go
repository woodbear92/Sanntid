package main
import(
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
//A backup should act similar as master when it states is set to master
func main(){
	elev := os.Args[1:][0]
	master:= os.Args[1:][1]
	elevNum, err := strconv.Atoi(elev)
	if err != nil || elevNum < 1{
		fmt.Println("Input error -- Make sure your elevator number is a positive integer")
		os.Exit(1)
	}
	masterNum, err := strconv.Atoi(master)
	if err != nil || elevNum < 1{
		fmt.Println("Input error")
		os.Exit(1)
	}
	InitBackUp(elevNum,masterNum)
}


func InitBackUp(elevNum int,mast int){
	fmt.Println("Elevator is:", elevNum)
	fmt.Println("Master is", mast)
  AllElevators   := [_numElevators] typedef.ElevatorData{}
	var masterWaitTimer *time.Timer = time.NewTimer(config.BackupWaitTimeMs*time.Millisecond)
	masterWaitTimer.Stop()
	masterData := typedef.MasterData{}
	masterData.State = typedef.MasterState_Backup
	masterData.Priorities[elevNum-1] = _numElevators

	master_elev := mast


	//Channels
	RXBackupElevatorData 			:= make(chan [_numElevators] typedef.ElevatorData)
	RXBackupBackupData				:= make(chan typedef.MasterData)
	TXBackupPing							:= make(chan int)
	MasterPingRx 						:= make(chan int)

	go bcast.Receiver(config.MasterPingPort, MasterPingRx)
	go bcast.Receiver(config.BackupElevatorDataTxPort, RXBackupElevatorData)
	go bcast.Receiver(config.BackupDataTxPort, RXBackupBackupData)
	go bcast.Transmitter(config.BackupPingPort+elevNum-1,TXBackupPing)
	//Send out a signal every .1s
	go func(){
		for {
			TXBackupPing <- elevNum
			time.Sleep(100*time.Millisecond)
		}
	}()

	//Print states
	go func(){
		for{
			if(masterData.State == typedef.MasterState_Backup){
				fmt.Println("I am backup")
				typedef.AllElevatorsPrint(AllElevators)
				typedef.PrintMasterData(masterData)
				time.Sleep(2000*time.Millisecond)

			}else {
				fmt.Println("I am master")
				typedef.AllElevatorsPrint(AllElevators)
				time.Sleep(2000*time.Millisecond)
			}
		}
	}()
	for{
		select {
		case allElev := <- RXBackupElevatorData:
			AllElevators = allElev
			masterWaitTimer.Reset(time.Duration(masterData.Priorities[elevNum-1])*config.BackupWaitTimeMs*time.Millisecond)
		case backupData := <- RXBackupBackupData:
			masterWaitTimer.Reset(time.Duration(masterData.Priorities[elevNum-1])*config.BackupWaitTimeMs*time.Millisecond)
			masterData.Priorities = backupData.Priorities
		case m := <-MasterPingRx:
			master_elev = m
			masterWaitTimer.Reset(time.Duration(masterData.Priorities[elevNum-1])*config.BackupWaitTimeMs*time.Millisecond)
		case <- masterWaitTimer.C:
				if (masterData.Priorities[elevNum-1] == 1){
					masterData.State = typedef.MasterState_Master
					masterData.Priorities[master_elev-1] = highestPriority
					masterData.Priorities = master.RecalculatePriorities(masterData.Priorities)
					go becomeMaster(elevNum,AllElevators,masterData)
				}
		}
		}
}

func becomeMaster(elevNum int,all_elev [_numElevators]typedef.ElevatorData, masterData typedef.MasterData){
	AllElevators := all_elev

	masterData.State = typedef.MasterState_Master
	masterData.Priorities[elevNum-1]=-1
	//Channels
	TXElevators     			:= make(chan [_numElevators] typedef.ElevatorData)
	RXElevators     			:= make(chan typedef.ElevatorData)
	TXBackupElevatorData 							:= make(chan [_numElevators] typedef.ElevatorData)
	TXBackupBackupData							:= make(chan typedef.MasterData)
	MasterPingTx 					:= make(chan int)
	elevatorDisconnectedChan  := make(chan int)
	backupDisconnectedChan  := make(chan int)
	backupConnectedChan  := make(chan int)


	//Create goroutines to check for elevator data received and check if backups are connected
	for e := 0; e < _numElevators; e++{
		go network.ReceiveDataFromElevators(RXElevators,e,elevatorDisconnectedChan)
		if (e != elevNum-1){
			go network.ReceiveBackupConfirmation(e,backupDisconnectedChan,backupConnectedChan)
			}
	}
	go bcast.Transmitter(config.MasterTXPort, TXElevators)
	go bcast.Transmitter(config.MasterPingPort, MasterPingTx)
	go bcast.Transmitter(config.BackupElevatorDataTxPort, TXBackupElevatorData)
	go bcast.Transmitter(config.BackupDataTxPort, TXBackupBackupData)
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
			TXBackupBackupData <-masterData
			TXBackupElevatorData <- AllElevators
			time.Sleep(100*time.Millisecond)
		}
	}()

	for{
		select{
		case e := <-RXElevators:
			AllElevators[e.Number] = e
			AllElevators = master.DistributeOrders(AllElevators, e)
			TXElevators <- AllElevators
		case e := <- elevatorDisconnectedChan:
			AllElevators[e].Connected = false
			AllElevators = master.RedistributeOrders(AllElevators, AllElevators[e])
			TXElevators <- AllElevators
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
