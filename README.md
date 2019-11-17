Elevator Project 
================

This is the repository for the "elevator project" in TTK4145 - Real-time Programming course.
The program language used is GO.

### How to run.
To run the executable write `./main elevator_number port_number`. The elevator number should be unique and starting at 1. Example: `./main 1 15657`. If the executable does not work, you can always try to rebuild the executable using the command `go build main.go`, or simply run the program with the command `go run main.go elevator_number port_number`.

Important! Elevator Number and port for network used must be set by the user.

### Modules

All modules (go packages) are developed in the group except the part of the network module and elevator I/O modules.

These modules can be found here:
* [elevator_control/elevio](https://github.com/TTK4145/driver-go)
* [network](https://github.com/TTK4145/Network-go)
 

Constants are found in config package and other types used throughout is found in typedef package

The logic for controlling the elevator is split up in Single and multiple elevator mode. Both found in logic folder.

The Elevator will run two different files depending if it is the master or one of the backups. The logic for the master is found in master folder. 

Cost function and the functions used to spawn master/backup processes are in the queue module.

Both the cost function and state machine algorithms/code is to a large degree inspired by the hand out material found at:

* [Single elevator algorithm](https://github.com/TTK4145/Project-resources/tree/master/elev_algo)
* [Cost functions](https://github.com/TTK4145/Project-resources/tree/master/cost_fns)
