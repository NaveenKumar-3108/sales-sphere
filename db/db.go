package db

import (
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Structure to hold database connection details
type DatabaseType struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DB       string
}

// structure to hold all db connection details used in this program
type AllUsedDatabases struct {
	Mysql DatabaseType
}

func LocalGORMDbConnect(pDBtype string) (*gorm.DB, error) {
	lDbDetails := new(AllUsedDatabases)
	lDbDetails.Init()

	lConnString := ""
	lLocalDBtype := ""

	var lErr error
	var lDataBaseConnection DatabaseType
	var lDBCon *gorm.DB
	var lDialector gorm.Dialector

	// get connection details
	if pDBtype == lDbDetails.Mysql.DB {
		lDataBaseConnection = lDbDetails.Mysql
		lLocalDBtype = lDbDetails.Mysql.DBType
	}
	log.Println(lDbDetails.Mysql, "qqmysql1", pDBtype)
	// Prepare connection string
	if lLocalDBtype == "mysql" {
		lConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", lDataBaseConnection.User, lDataBaseConnection.Password, lDataBaseConnection.Server, lDataBaseConnection.Port, lDataBaseConnection.Database)
		lDialector = mysql.Open(lConnString)
	} else {
		log.Println("DB details not found")
	}

	log.Println("mysql", lConnString, lLocalDBtype)
	lDBCon, lErr = gorm.Open(lDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if lErr != nil {
		// pDebug.Log(helpers.Elog, "LocalGORMDbConnect 001:Invalid DB Details", lErr)
		log.Println("Invalid DB Details")
		return lDBCon, fmt.Errorf("Invalid DB Details")
	}

	return lDBCon, lErr
}
