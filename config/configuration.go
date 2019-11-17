package config

import (
	"time"
)

//Hardware configuration:
const NumFloors int = 4
const NumElevators int = 3
const NumOrderButtons int = 3

//Time for timer
const DoorOpenTimeMs time.Duration = 3000
const ReconnectionTimeMs time.Duration = 5000
const TravelTimeMs time.Duration = 3500
const BackupWaitTimeMs time.Duration = 5000

//Ports used
const ElevatorTXPort int = 11452
const MasterTXPort int = 16001
const BackupPingPort int = 15600
const MasterPingPort int = 16835
const BackupDataTxPort int = 15001
const BackupElevatorDataTxPort int = 15002
