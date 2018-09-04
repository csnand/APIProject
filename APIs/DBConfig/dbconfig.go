package DBConfig

import (
	"database/sql"
	"fmt"

	//only init() in pq package are needed here, no other explicit function calls
	_ "github.com/lib/pq"
)

var DBCONN *sql.DB

const (
	dbhost = "localhost"
	dbport = "5432"
	dbuser = "dev"
	dbpass = "dev"
	dbname = "api_development"
)

//InitDb initialise db info and establish db connection
func InitDb() error {

	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config["host"], config["port"],
		config["user"], config["pass"], config["name"])

	DBCONN, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	err = DBCONN.Ping()
	if err != nil {
		return err
	}
	return nil
}

func PInitDb() (*sql.DB, error) {

	var DBCONN *sql.DB
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config["host"], config["port"],
		config["user"], config["pass"], config["name"])

	DBCONN, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = DBCONN.Ping()
	if err != nil {
		return nil, err
	}
	return DBCONN, nil
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	conf["host"] = dbhost
	conf["port"] = dbport
	conf["user"] = dbuser
	conf["pass"] = dbpass
	conf["name"] = dbname

	// conf["host"] = os.Getenv("DBHOST")
	// conf["port"] = os.Getenv("DBPORT")
	// conf["user"] = os.Getenv("DBUSER")
	// conf["pass"] = os.Getenv("DBPASS")
	// conf["name"] = os.Getenv("DBNAME")

	return conf
}
