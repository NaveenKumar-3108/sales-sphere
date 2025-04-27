package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sales-sphere/apps/refresh"
	"sales-sphere/apps/revenue"
	"sales-sphere/apps/scheduler"
	"sales-sphere/db"

	"github.com/gorilla/mux"
)

type Product struct {
	ProductID    int    `gorm:"primaryKey"`
	ProductName  string `gorm:"not null"`
	CategoryName string
	Description  string
	UnitPrice    float64 `gorm:"not null"`
}
type User struct {
	UserID   int    `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	Category string `gorm:"column:category"`
}

func main() {
	log.Println("Server Started")

	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	lErr := db.BuildConnection()
	if lErr != nil {
		log.Println("Error in DB conection")
	}
	go scheduler.LoadSalesData()

	router := CreateRouter()

	log.Fatal(http.ListenAndServe(":29022", router))
}

func CreateRouter() *mux.Router {
	Router := mux.NewRouter()
	Router.HandleFunc("/api/refresh", refresh.RefreshSalesData).Methods("GET")
	Router.HandleFunc("/api/getRevenue", revenue.GetRevenue).Methods("GET")

	return Router
}
