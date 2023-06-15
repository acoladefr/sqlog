package sqlog

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Logger struct {
	table    string
	allowSQL bool
}

type Opts func(*Logger)

func WithTable(table string) Opts {
	return func(l *Logger) {
		l.table = table
	}
}

func WithSQL(allowSQL bool) Opts {
	return func(l *Logger) {
		l.allowSQL = allowSQL
	}
}

func NewLogger(opts ...Opts) *Logger {
	l := &Logger{}
	for _, option := range opts {
		option(l)
	}
	return l
}

func (l Logger) Table() string {
	return l.table
}

func (l *Logger) Init() error {
	l.Log(DEBUG, CONSOLE, "sqlog-init", "trying to connect to db")
	db, err := sql.Open("mysql", os.Getenv("MYSQL_USERNAME")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB"))

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + l.Table() + ` (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		timestamp BIGINT,
		origin VARCHAR(255),
		type TINYINT,
		message VARCHAR(255),
		system_involved VARCHAR(255)
	);`)

	if err != nil {
		return err
	}

	l.Log(DEBUG, CONSOLE, "sqlog-init", "db ready; sqlog ready")
	return nil
}

type LEVEL uint8

type TARGET uint8

const (
	FATAL LEVEL = iota + 1
	ERROR
	WARNING
	INFO
	DEBUG
)

const (
	CONSOLE TARGET = iota
	DATABASE
	BOTH
)

func (l LEVEL) getLabel() string {
	switch l {
	case FATAL:
		return "\033[31mfatal\033[0m"
	case ERROR:
		return "\033[31merror\033[0m"
	case WARNING:
		return "\033[33mwarning\033[0m"
	case INFO:
		return "\033[32minfo\033[0m"
	case DEBUG:
		return "\033[34mdebug\033[0m"
	}
	return ""
}

func (l *Logger) Log(level LEVEL, target TARGET, system string, format string, args ...interface{}) {
	message := format
	for i, arg := range args {
		placeholder := "{" + fmt.Sprint(i) + "}"
		message = strings.Replace(message, placeholder, fmt.Sprint(arg), -1)
	}
	switch target {
	case CONSOLE:
		date, _ := exec.Command("bash", "-c", "echo -n $(date +'%Y-%m-%d %H:%M:%S')").Output()
		fmt.Fprintf(os.Stdout, "("+string(date)+")["+level.getLabel()+"; \033[36m"+system+"\033[0m]: "+message+"\n")
	case DATABASE:
		if l.allowSQL && l.table != "" {
			db, err := sql.Open("mysql", os.Getenv("MYSQL_USERNAME")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB"))

			if err != nil {
				l.Log(ERROR, CONSOLE, "sqlog-db-connection", "failed to open connection to db: {0}", err)
			}

			stmt, err := db.Prepare("INSERT INTO " + l.Table() + " (timestamp, origin, type, message, system_involved) VALUES (?, ?, ?, ?, ?)")

			if err != nil {
				l.Log(ERROR, CONSOLE, "sqlog-prepare-stmt", "failed to prepare statement: {0}", err)
			}

			defer stmt.Close()

			t, _ := exec.Command("bash", "-c", "echo -n $(date +%s)").Output()
			timestamp, _ := strconv.ParseInt(string(t), 10, 64)

			o, _ := exec.Command("bash", "-c", "echo -n $(hostname)").Output()
			origin := string(o)

			_, err = stmt.Exec(timestamp, origin, level, message, system)
			if err != nil {
				l.Log(ERROR, CONSOLE, "sqlog-execute-stmt", "failed to execute statement: {0}", err)
			}
		}
	case BOTH:
		l.Log(level, CONSOLE, system, message)
		l.Log(level, DATABASE, system, message)
	}
}
