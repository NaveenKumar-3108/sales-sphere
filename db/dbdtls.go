package db

import (
	"fmt"
	"sales-sphere/common"
	"strconv"
)

const (
	MYSQL = "mysql"
)

// Initializing DB Details
func (d *AllUsedDatabases) Init() {
	dbconfig := common.ReadTomlConfig("./dbconfig.toml")

	d.Mysql.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlServer"])
	d.Mysql.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlPort"]))
	d.Mysql.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlUser"])
	d.Mysql.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlPassword"])
	d.Mysql.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlDatabase"])
	d.Mysql.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MysqlDBType"])
	d.Mysql.DB = MYSQL

}
