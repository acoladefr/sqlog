package main

import "acolade.ch/sqlog"

var logger *sqlog.Logger

func main() {
	logger = sqlog.NewLogger(sqlog.WithSQL(true), sqlog.WithTable("logs"))

	_, err := logger.Init()
	if err != nil {
		panic(err)
	}

	logger.Log(sqlog.INFO, sqlog.BOTH, "info-log-in-database-and-console", "Info log")
	logger.Log(sqlog.WARNING, sqlog.CONSOLE, "warning-log-in-console", "warning log")
	logger.Log(sqlog.ERROR, sqlog.DATABASE, "warning-log-in-database", "error log")
	logger.Log(sqlog.DEBUG, sqlog.CONSOLE, "debug-log-in-console", "debug log")
	logger.Log(sqlog.FATAL, sqlog.BOTH, "fatal-log-in-console", "fatal log")
}
