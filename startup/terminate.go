package startup

import (
	"os"
	"os/signal"
	"syscall"
)

func Terminate(log Logger) {
	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, syscall.SIGINT, syscall.SIGTERM)
	<-sCh
	log.Infof("terminating")
	// k8s is in charge now so undo handling of signals
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)
}
