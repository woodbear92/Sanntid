package queue

import (
	"fmt"
	"os/exec"
	"strconv"
)

func Spawn(rank string, elev int, master int) bool {
	fmt.Println("Spawning...")
	if rank == "master" {
		spawnNewMaster(elev)
		return true
	} else if rank == "backup" {
		spawnNewBackup(elev, master)
		return true
	} else {
		return false
	}
}

func spawnNewBackup(elev int, master int) {
	str := "go run backup.go " + strconv.Itoa(elev) + " " + strconv.Itoa(master)
	(exec.Command("gnome-terminal", "-x", "sh", "-c", str)).Run()
	fmt.Println("Backup spawning")
}

func spawnNewMaster(elev int) {
	str := "go run master.go " + strconv.Itoa(elev)
	(exec.Command("gnome-terminal", "-x", "sh", "-c", str)).Run()
	fmt.Println("Master spawning")
}
