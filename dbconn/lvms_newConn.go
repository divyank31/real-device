package dbconn

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var LVMSConn *sql.DB

type lvmsDbConfig struct {
	Driver       string
	Protocol     string
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func getlvmsDbConfig() (dbConfig lvmsDbConfig, err error) {
	// dbHost := os.Getenv("ML_DB_HOST")
	// dbPort := os.Getenv("ML_DB_PORT")
	// dbUsername := os.Getenv("DB_USERNAME")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("LVMS_DB_NAME")

	dbConfig.Host = "localhost"               //dbHost
	dbConfig.Port = "3306"                    //dbPort
	dbConfig.Username = "root"                //dbUsername
	dbConfig.Password = "Qwerty1!"            //dbPassword
	dbConfig.DatabaseName = "lambda_lvms_new" //dbName
	return dbConfig, nil
}

func InitLVMSConn() {

	configuration, err := getlvmsDbConfig()

	if err != nil {
		fmt.Println("fetch-dbconfig-connection-failed : ", err.Error())
	}

	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configuration.Username, configuration.Password, configuration.Host, configuration.Port, configuration.DatabaseName)
	LVMSConn, err = sql.Open("mysql", dbString)
	fmt.Println("lvms connection " + dbString)
	if err != nil {
		panic(err.Error())
	}
	//Pool.SetMaxIdleConns(c)
	//Pool.SetMaxOpenConns(constants.MYSQL_MAX_OPEN_CONNECTION)
	//Pool.SetConnMaxLifetime(constants.MYSQL_MAX_CONNECTION_LIFETIME)

	err = LVMSConn.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		fmt.Println("Error on opening database connection: ", err.Error())
	}
}

func GetLVMSConn() *sql.DB {
	return LVMSConn
}
