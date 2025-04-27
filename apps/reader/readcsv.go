package reader

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type Record struct {
	OrderID       int
	ProductID     int
	CustomerID    int
	ProductName   string
	Category      string
	Region        string
	DateOfSale    string
	Quantity      int
	UnitPrice     float64
	Discount      float64
	ShippingCost  float64
	PaymentMethod string
	CustomerName  string
	CustomerEmail string
	CustomerAddr  string
}

func ReadCSV(filePath string) ([]Record, error) {
	log.Println("ReadCsv(+)")

	var lRecordList []Record
	f, lErr := os.Open(filePath)
	if lErr != nil {
		return nil, lErr
	}
	defer f.Close()

	reader := csv.NewReader(f)

	lRecords, lErr := reader.ReadAll()

	if lErr != nil {
		return lRecordList, lErr
	}
	for i, lRow := range lRecords {

		if i == 0 {
			continue // skip header
		}

		lOrderID, _ := strconv.Atoi(lRow[0])
		lProductID, _ := strconv.Atoi(lRow[1])
		lCustomerID, _ := strconv.Atoi(lRow[2])
		lQuantity, _ := strconv.Atoi(lRow[7])

		lUnitPrice, _ := strconv.ParseFloat(lRow[8], 64)
		lDiscount, _ := strconv.ParseFloat(lRow[9], 64)
		lShippingCost, _ := strconv.ParseFloat(lRow[10], 64)

		lRecordList = append(lRecordList, Record{
			OrderID:       lOrderID,
			ProductID:     lProductID,
			CustomerID:    lCustomerID,
			ProductName:   lRow[3],
			Category:      lRow[4],
			Region:        lRow[5],
			DateOfSale:    lRow[6],
			Quantity:      lQuantity,
			UnitPrice:     lUnitPrice,
			Discount:      lDiscount,
			ShippingCost:  lShippingCost,
			PaymentMethod: lRow[11],
			CustomerName:  lRow[12],
			CustomerEmail: lRow[13],
			CustomerAddr:  lRow[14],
		})
	}
	log.Println("ReadCsv(+)")
	return lRecordList, nil
}
