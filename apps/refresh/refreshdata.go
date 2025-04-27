package refresh

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sales-sphere/apps/reader"
	"sales-sphere/common"
	"sales-sphere/db"

	"strings"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ProductID   int     `gorm:"column:ProductID;primaryKey;autoIncrement"`
	ProductName string  `gorm:"column:ProductName"`
	Category    string  `gorm:"column:Category;type:varchar(100)"`
	Description string  `gorm:"column:Description;type:text"`
	UnitPrice   float64 `gorm:"column:UnitPrice"`
}

type Customer struct {
	CustomerID    int    `gorm:"column:CustomerID;primaryKey;autoIncrement"`
	CustomerName  string `gorm:"column:CustomerName"`
	CustomerEmail string `gorm:"column:CustomerEmail;not null"`
	CustomerAddr  string `gorm:"column:CustomerAddress;type:text"`
	Region        string `gorm:"column:Region;type:varchar(100)"`
}

type Order struct {
	OrderID       int       `gorm:"column:OrderID;primaryKey"`
	CustomerID    int       `gorm:"column:CustomerID;not null"`
	DateOfSale    time.Time `gorm:"column:DateOfSale;not null"`
	PaymentMethod string    `gorm:"column:PaymentMethod;type:varchar(50)"`
	ShippingCost  float64   `gorm:"column:ShippingCost;type:decimal(10,2)"`
	Discount      float64   `gorm:"column:Discount;type:decimal(5,2)"`
	Customer      Customer  `gorm:"foreignKey:CustomerID"`
}

type OrderItems struct {
	OrderItemID  int     `gorm:"column:OrderItemID;primaryKey;autoIncrement"`
	OrderID      int     `gorm:"foreignKey:column:OrderID;not null"`
	ProductID    int     `gorm:"column:ProductID;not null"`
	QuantitySold int     `gorm:"column:QuantitySold;not null"`
	UnitPrice    float64 `gorm:"column:UnitPrice;not null"`
	Order        Order   `gorm:"foreignKey:OrderID;references:OrderID"`
	Product      Product `gorm:"foreignKey:ProductID;references:ProductID"`
}

type Response struct {
	Status string `json:"status"`
	ErrMsg string `json:"errmsg"`
}

func RefreshSalesData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	log.Println("RefreshSalesData(+)")
	if strings.EqualFold("GET", r.Method) {
		var lResp Response
		lResp.Status = common.SUCCESS
		lTomlConfig := common.ReadTomlConfig("./toml/config.toml")

		lPath := fmt.Sprintf("%v", lTomlConfig.(map[string]interface{})["path"])
		lSalesList, lErr := reader.ReadCSV(lPath)
		if lErr != nil {
			log.Println("Error:RLS01", lErr)
			lResp.Status = common.ERROR
			lResp.ErrMsg = lErr.Error()
			goto Marshal
		} else {
			for _, lRecord := range lSalesList {

				lErr := ProcessRecords(lRecord)

				if lErr != nil {
					log.Println("Error:RLS02", lErr)
					lResp.Status = common.ERROR
					lResp.ErrMsg = lErr.Error()
					goto Marshal
				}
			}
			lResp.ErrMsg = "Data loaded successfully"
		}

	Marshal:
		lData, lErr := json.Marshal(lResp)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			fmt.Fprint(w, string(lData))
		}
	}
	log.Println("RefreshSalesData(-)")
}

func ProcessRecords(pRecord reader.Record) error {
	log.Println("ProcessRecords(+)")

	lProduct := Product{

		ProductName: pRecord.ProductName,
		Category:    pRecord.Category,
		UnitPrice:   pRecord.UnitPrice,
	}
	if lErr := InsertProduct(&lProduct); lErr != nil {
		return lErr
	}

	// Insert or update the customer
	lCustomer := Customer{

		CustomerName:  pRecord.CustomerName,
		CustomerEmail: pRecord.CustomerEmail,
		CustomerAddr:  pRecord.CustomerAddr,
		Region:        pRecord.Region,
	}
	if lErr := InsertCustomer(&lCustomer); lErr != nil {
		return lErr
	}

	lFormat := "2006-01-02"

	lParsedTime, lErr := time.Parse(lFormat, pRecord.DateOfSale)
	if lErr != nil {
		log.Println("Error parsing date:", lErr)
		return lErr
	}

	lOrder := Order{
		OrderID:       pRecord.OrderID,
		CustomerID:    lCustomer.CustomerID,
		DateOfSale:    lParsedTime,
		PaymentMethod: pRecord.PaymentMethod,
		ShippingCost:  pRecord.ShippingCost,
		Discount:      pRecord.Discount,
	}
	if lErr := InsertOrder(&lOrder); lErr != nil {
		return lErr
	}

	// Insert the order item (assuming OrderID and ProductID make it unique)
	lOrderItem := OrderItems{
		OrderID:      lOrder.OrderID,
		ProductID:    lProduct.ProductID,
		QuantitySold: pRecord.Quantity,
		UnitPrice:    pRecord.UnitPrice,
	}
	if lErr := InsertOrderItem(lOrderItem); lErr != nil {
		return lErr
	}

	log.Println("ProcessRecords(-)")
	return nil
}

func InsertCustomer(customer *Customer) error {
	log.Println("InsertCustomer(+)")
	lResult := db.GormDBConnection.Where("CustomerID = ?", customer.CustomerEmail).First(&customer)
	if lResult.Error != nil && lResult.Error != gorm.ErrRecordNotFound {
		log.Println("Error:RLIC01", lResult.Error)
		return lResult.Error
	}

	if lResult.RowsAffected == 0 { //no record fount insert done
		if lErr := db.GormDBConnection.Create(&customer).Error; lErr != nil {
			log.Println("Error:RLIC02", lErr)
			return lErr
		}
	} else { //if record alredy found  update done

		if lErr := db.GormDBConnection.Save(&customer).Error; lErr != nil {
			log.Println("Error:RLIC03", lResult.Error)
			return lErr
		}
	}
	log.Println("InsertCustomer(-)")
	return nil
}

func InsertProduct(pProduct *Product) error {
	log.Println("InsertProduct(+)")

	lResult := db.GormDBConnection.Where("ProductID = ?", pProduct.ProductID).First(&pProduct)
	if lResult.Error != nil && lResult.Error != gorm.ErrRecordNotFound {
		log.Println("Error:RLIP01", lResult.Error)

		return lResult.Error
	}

	if lResult.RowsAffected == 0 {
		if lErr := db.GormDBConnection.Create(&pProduct).Error; lErr != nil {
			log.Println("Error:RLIP02", lErr)
			return lErr
		}
	} else {

		if lErr := db.GormDBConnection.Save(&pProduct).Error; lErr != nil {
			log.Println("Error:RLIP03", lErr)
			return lErr
		}
	}
	log.Println("InsertProduct(-)")
	return nil
}

func InsertOrder(pOrder *Order) error {
	log.Println("InsertOrder(+)")
	lResult := db.GormDBConnection.Where("OrderID = ?", pOrder.OrderID).First(&pOrder)
	if lResult.Error != nil && lResult.Error != gorm.ErrRecordNotFound {
		log.Println("Error:RLI001", lResult.Error)
		return lResult.Error
	}

	if lResult.RowsAffected == 0 {
		if lErr := db.GormDBConnection.Create(&pOrder).Error; lErr != nil {
			log.Println("Error:RLI002", lErr)
			return lErr
		}
	} else {

		if lErr := db.GormDBConnection.Save(&pOrder).Error; lErr != nil {
			log.Println("Error:RLI003", lErr)
			return lErr
		}
	}
	log.Println("InsertOrder(-)")
	return nil
}

func InsertOrderItem(pOrderItem OrderItems) error {
	log.Println("InsertOrderItem(+)")
	if lErr := db.GormDBConnection.Create(&pOrderItem).Error; lErr != nil {
		log.Println("Error:RLOI01", lErr)
		return lErr
	}
	log.Println("InsertOrderItem(-)")
	return nil
}
