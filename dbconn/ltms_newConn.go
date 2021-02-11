package dbconn

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var LTMSConn *sql.DB

type ltmsDbConfig struct {
	Driver       string
	Protocol     string
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func getltmsDbConfig() (dbConfig ltmsDbConfig, err error) {
	// dbHost := os.Getenv("ML_DB_HOST")
	// dbPort := os.Getenv("ML_DB_PORT")
	// dbUsername := os.Getenv("DB_USERNAME")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("LTMS_DB_NAME")

	dbConfig.Host = "localhost"               //dbHost
	dbConfig.Port = "3306"                    //dbPort
	dbConfig.Username = "root"                //dbUsername
	dbConfig.Password = "Qwerty1!"            //dbPassword
	dbConfig.DatabaseName = "lambda_ltms_new" //dbName
	return dbConfig, nil
}

func InitLTMSConn() {

	configuration, err := getltmsDbConfig()

	if err != nil {
		fmt.Println("fetch-dbconfig-connection-failed : ", err.Error())
	}

	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configuration.Username, configuration.Password, configuration.Host, configuration.Port, configuration.DatabaseName)
	fmt.Println("ltms connection " + dbString)
	LTMSConn, err = sql.Open("mysql", dbString)

	if err != nil {
		panic(err.Error())
	}
	//Pool.SetMaxIdleConns(c)
	//Pool.SetMaxOpenConns(constants.MYSQL_MAX_OPEN_CONNECTION)
	//Pool.SetConnMaxLifetime(constants.MYSQL_MAX_CONNECTION_LIFETIME)

	err = LTMSConn.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		fmt.Println("Error on opening database connection: ", err.Error())
	}
}

func GetLTMSConn() *sql.DB {
	return LTMSConn
}
