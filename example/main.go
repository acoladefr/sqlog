package main

import "acolade.ch/sqlog"

var logger *sqlog.Logger

func main() {
	logger = sqlog.NewLogger(sqlog.WithSQL(true), sqlog.WithTable("logs"))

	_, err := logger.Init()
	if err != nil {
		panic(err)
	}

	logger.Log(sqlog.INFO, sqlog.BOTH, "test-log", "Hello World!")
}
