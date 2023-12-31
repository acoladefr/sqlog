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

import  "github.com/acoladech/sqlog" 

var  logger  *sqlog.Logger

func  main()  {
    logger  = sqlog.NewLogger(sqlog.WithSQL(true), sqlog.WithTable("logs"))

    //logger can insert log in the table "logs" of the database

    _,  err  := logger.Init()

    if err !=  nil  {
	    panic(err)
    }

    // the Init() method will create the desired table if it does not exist
    // and return an error if the table cannot be initialized
    // NOTE: if you do not use sql, you do not have to call Init()

    address := "127.0.0.1"
    port    := 45912

    logger.Log(sqlog.INFO, sqlog.BOTH,  "daemon-initialization",  "daemon listening on {0}:{1}", address, port)
    // log in both database and console of type info that
    // gives information about the "daemon-initialization" process
    // and says that the daemon is listening on port 45912
}
```
