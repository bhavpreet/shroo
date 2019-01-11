package main

import "github.com/sirupsen/logrus"

func Init(interface{}) interface{} {
	return nil
}

func Print(a interface{}) interface{} {
	go func() {
		logrus.Info(a)
	}()
	return nil
}
