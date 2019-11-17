package network

import (
	"../config"
	"../typedef"
	"./bcast"
	"time"
)

func ReceiveDataFromElevators(receive_chan chan<- typedef.ElevatorData, elev int, disconnected_chan chan<- int) {
	checkConnectionTimer := time.NewTimer(config.ReconnectionTimeMs * time.Millisecond)
	elevatorRX := make(chan typedef.ElevatorData)

	go bcast.Receiver(config.ElevatorTXPort+elev, elevatorRX)

	for {
		select {
		case e := <-elevatorRX:
			checkConnectionTimer.Stop()
			e.Connected = true
			if e.HardwareFailure == true {
				//Same behavior if hardware failure accures as if it is disconnected
				e.Connected = false
				disconnected_chan <- elev
			}
			receive_chan <- e
			checkConnectionTimer.Reset(config.ReconnectionTimeMs * time.Millisecond)
		case <-checkConnectionTimer.C:
			disconnected_chan <- elev
			checkConnectionTimer.Stop()
		}
	}
}

func ReceiveBackupConfirmation(elev int, disconnect_backup chan<- int, connected_backup chan<- int) {
	checkConnectionTimer := time.NewTimer(config.ReconnectionTimeMs * time.Millisecond)
	RxBackup := make(chan int)

	go bcast.Receiver(config.BackupPingPort+elev, RxBackup)

	for {
		select {
		case <-RxBackup:
			connected_backup <- elev
			checkConnectionTimer.Reset(config.ReconnectionTimeMs * time.Millisecond)
		case <-checkConnectionTimer.C:
			disconnect_backup <- elev
		}
	}
}
