package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func CtrlC() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	<-sigCh
	fmt.Println("\rCtrl+C is not responding")
	<-sigCh
	os.Exit(0)
}

func CtrlCKill() {
	var pid int
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	<-sigCh
	fmt.Println("\rEnter the PID of the process you want to kill: ")
	_, err := fmt.Scan(&pid)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	err = syscall.Kill(pid, syscall.SIGINT)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Printf("Process %v killed: ", pid)
}
