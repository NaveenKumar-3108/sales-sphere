package revenue

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sales-sphere/common"
	"sales-sphere/db"
	"strconv"
	"strings"
	"time"
)

type RevenueResp struct {
	TotalRevenue float64 `json:"total_revenue"`
	Status       string  `json:"status"`
	Errmsg       string  `json:"errmsg"`
}

func GetRevenue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "StartDate,EndDate,ProductID,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("GetRevenue(+)")
	if strings.EqualFold("GET", r.Method) {
		var lResp RevenueResp
		lResp.Status = common.SUCCESS

		lStartDate := r.Header.Get("StartDate")
		lEndDate := r.Header.Get("EndDate")
		lProductID := r.Header.Get("ProductID")
		if lStartDate != "" || lEndDate != "" {
			if lProductID == "" { //productid is empty provide total revenue
				lToatalRevenue, lErr := GetTotalRevenue(lStartDate, lEndDate)

				if lErr != nil {
					log.Println("Error:RGR01", lErr)
					lResp.Status = common.ERROR
					lResp.Errmsg = lErr.Error()
					goto Marshal
				} else {
					lResp.TotalRevenue = lToatalRevenue
				}
			} else { // product is present provide product based revenue
				lToatalRevenue, lErr := GetRevenueByProduct(lStartDate, lEndDate, lProductID)

				if lErr != nil {
					log.Println("Error:RGR02", lErr)
					lResp.Status = common.ERROR
					lResp.Errmsg = lErr.Error()
					goto Marshal
				} else {
					lResp.TotalRevenue = lToatalRevenue
				}
			}
		} else {
			lResp.Status = common.ERROR
			lResp.Errmsg = "Invalid date"
		}

	Marshal:
		lData, lErr := json.Marshal(lResp)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			fmt.Fprint(w, string(lData))
		}
	}
	log.Println("GetRevenue(-)")
}
func GetTotalRevenue(pStartDate, pEndDate string) (float64, error) {
	log.Println("GetTotalRevenue(+)")
	var lTotalRevenue float64
	lFormat := "2006-01-02"
	lParsedStartTime, lErr := time.Parse(lFormat, pStartDate)
	if lErr != nil {
		log.Println("RGRT01:", lErr)
		return lTotalRevenue, lErr
	}
	lParsedEndTime, lErr := time.Parse(lFormat, pEndDate)
	if lErr != nil {
		log.Println("Error :RGRT02", lErr)
		return lTotalRevenue, lErr
	}
	err := db.GormDBConnection.Table("order_items").
		Select("SUM(order_items.QuantitySold * order_items.UnitPrice) AS total_revenue").
		Joins("INNER JOIN orders ON orders.OrderID = order_items.OrderID").
		Where("orders.DateOfSale BETWEEN ? AND ?", lParsedStartTime, lParsedEndTime).
		Scan(&lTotalRevenue).Error
	if err != nil {
		return lTotalRevenue, err
	}
	log.Println("GetTotalRevenue(-)")
	return lTotalRevenue, nil
}

func GetRevenueByProduct(pStartDate, pEndDate, pProductID string) (float64, error) {
	log.Println("GetRevenueByProduct(+)")
	var lProductRevenue float64
	lFormat := "2006-01-02"
	lParsedStartTime, lErr := time.Parse(lFormat, pStartDate)
	if lErr != nil {
		log.Println("Error:RGRP01", lErr)
		return lProductRevenue, lErr
	}
	lParsedEndTime, lErr := time.Parse(lFormat, pEndDate)
	if lErr != nil {
		log.Println("Error:RGRP02", lErr)
		return lProductRevenue, lErr
	}
	lProductID, _ := strconv.Atoi(pProductID)

	// Query to get total revenue for the specified product ID
	err := db.GormDBConnection.Table("order_items").
		Select("SUM(order_items.QuantitySold * order_items.UnitPrice) AS product_revenue").
		Joins("INNER JOIN orders ON orders.OrderID = order_items.OrderID").
		Where("orders.DateOfSale BETWEEN ? AND ? AND order_items.ProductID = ?", lParsedStartTime, lParsedEndTime, lProductID).
		Scan(&lProductRevenue).Error
	if err != nil {
		return lProductRevenue, err
	}
	log.Println("GetRevenueByProduct(-)")
	return lProductRevenue, nil
}
