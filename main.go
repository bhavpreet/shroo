package main

import (
	"github.com/bhavpreet/shroo/controller"
	"github.com/sirupsen/logrus"
)

func runController() {
	controller.Init()
	controller.Run()
}
func main() {
	logrus.SetLevel(logrus.InfoLevel)
	runController()
}
