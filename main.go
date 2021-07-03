package main

import (
	"github.com/sirupsen/logrus"
)

func runController() {
	ControllerInit()
	ControllerRun()
}
func main() {
	logrus.SetLevel(logrus.InfoLevel)
	runController()
}
