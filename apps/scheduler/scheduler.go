package scheduler

import (
	"fmt"
	"log"
	"sales-sphere/apps/reader"
	"sales-sphere/apps/refresh"
	"sales-sphere/common"
	"strconv"
	"time"
)

func LoadSalesData() {
	log.Println("LoadSalesData(+)")
	lTomlConfig := common.ReadTomlConfig("./toml/config.toml")
	lTime := fmt.Sprintf("%v", lTomlConfig.(map[string]interface{})["Hours"])
	lHours, _ := strconv.Atoi(lTime)
	interval := time.Duration(lHours) * time.Hour
	nextScheduledTime := time.Now().Add(interval)

	// Log when the schedule starts with the next run time
	log.Printf("Schedule started: ProcessCsv will run at %v (after %d hours)", nextScheduledTime.Format("2006-01-02 03:04:05 PM"), lHours)

	time.AfterFunc(interval, func() { //procees csv a specfic interval based on cofig hours
		ProcessCsv()    //to process
		LoadSalesData() //to resechudle
	})
	log.Println("LoadSalesData(-)")
}
func ProcessCsv() {
	log.Println("ProcessCsv(+)")
	lTomlConfig := common.ReadTomlConfig("./toml/config.toml")
	lPath := fmt.Sprintf("%v", lTomlConfig.(map[string]interface{})["path"])
	lSalesList, lErr := reader.ReadCSV(lPath)
	if lErr != nil {
		log.Println("Error:SPS01", lErr)
		return
	} else {
		for _, lRecord := range lSalesList {

			lErr := refresh.ProcessRecords(lRecord)

			if lErr != nil {
				log.Println("Error:SPS02", lErr)
				return
			}
		}
	}
	log.Println("ProcessCsv(-)")
}
