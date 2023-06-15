# sqlog
A logging library written in **Go** that can trace logs in a **MySQL** database

## Requirements
sqlog will use **environment variables** to connect to the database. (only if you want to insert logs in a db)

### List of needed environment variables:

 - MYSQL_HOST (ip address of the server)
 - MYSQL_PORT (port that the database uses)
 - MYSQL_USERNAME (your username)
 - MYSQL_PASSWORD (your password)
 - MYSQL_DB (your database's name)

## Setup
1. Add the sqlog library to the app:
```bash
go get github.com/acoladech/sqlog
```
## Example
You can take a look at file **example/main.go** but here is a classic example.

```go
package  main 

import  "acolade.ch/sqlog" 

var  logger  *sqlog.Logger

func  main()  {
logger  = sqlog.NewLogger(sqlog.WithSQL(true), sqlog.WithTable("logs"))

_,  err  := logger.Init()

if err !=  nil  {
	panic(err)
}

// log in both database and console of type info that
// gives information about the "daemon-initialization" process
// and says that the daemon is listening on port 45912
logger.Log(sqlog.INFO, sqlog.BOTH,  "daemon-initialization",  "daemon listening on :45912")

}
```
