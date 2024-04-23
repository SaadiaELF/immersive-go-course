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
