package db

import (
	"log"

	"gorm.io/gorm"
)

var GormDBConnection *gorm.DB

func BuildConnection() error {
	var lErr error
	GormDBConnection, lErr = LocalGORMDbConnect(MYSQL)
	if lErr != nil {
		log.Println("Error in DB connect")
		return lErr
	}
	return nil
}
