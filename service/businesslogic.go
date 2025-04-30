package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"lumel/models"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Conn struct {
	db *gorm.DB
}
type Service interface {
	ReadDataFromCSV(filePath string) error
	GetAverageOrderValue(startDate, endDate time.Time) (float64, error)
	GetTotalOrders(startDate, endDate time.Time) (int, error)
	GetTotalCustomers(startDate, endDate time.Time) (int, error)
}

func (s *Conn) ReadDataFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err

	}
	defer file.Close()
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		orderID := record[0]
		productID := record[1]
		customerID := record[2]
		productName := record[3]
		category := record[4]
		region := record[5]
		dateOfSale, _ := time.Parse("2025-04-30", record[6])
		quantitySold := atoi(record[7])
		unitPrice := atof(record[8])
		discount := atof(record[9])
		shippingCost := atof(record[10])
		paymentMethod := record[11]
		customerName := record[12]
		customerEmail := record[13]
		customerAddress := record[14]

		var product models.Product
		s.db.FirstOrCreate(&product, models.Product{ID: productID})

		var customer models.Customer
		s.db.FirstOrCreate(&customer, models.Customer{ID: customerID})

		order := models.Order{
			OrderID:         orderID,
			CustomerID:      customerID,
			ProductID:       productID,
			QuantitySold:    quantitySold,
			UnitPrice:       unitPrice,
			Discount:        discount,
			ShippingCost:    shippingCost,
			PaymentMethod:   paymentMethod,
			DateOfSale:      dateOfSale,
			ProductName:     productName,
			Category:        category,
			Region:          region,
			CustomerName:    customerName,
			CustomerEmail:   customerEmail,
			CustomerAddress: customerAddress,
			Product:         product,
			Customer:        customer,
		}

		s.db.Create(&order)
	}

	fmt.Println("Data refreshed successfully from CSV.")
	return nil
}
func atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

func atof(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return val
}
func (s *Conn) GetTotalCustomers(startDate, endDate time.Time) (int, error) {
	var totalCustomers int64
	err := s.db.Model(&models.Customer{}).Joins("JOIN orders ON orders.customer_id = customers.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Count(&totalCustomers).Error
	if err != nil {
		return 0, fmt.Errorf("error fetching total customers: %v", err)
	}
	return int(totalCustomers), nil
}

func (s *Conn) GetTotalOrders(startDate, endDate time.Time) (int, error) {
	var totalOrders int64
	err := s.db.Model(&models.Order{}).Where("date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Count(&totalOrders).Error
	if err != nil {
		return 0, fmt.Errorf("error fetching total orders: %v", err)
	}
	return int(totalOrders), nil
}

func (s *Conn) GetAverageOrderValue(startDate, endDate time.Time) (float64, error) {
	var totalRevenue float64
	err := s.db.Model(&models.Order{}).Where("date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Select("SUM(quantity_sold * unit_price * (1 - discount))").Scan(&totalRevenue).Error
	if err != nil {
		return 0, fmt.Errorf("error fetching total revenue: %v", err)
	}

	totalOrders, err := s.GetTotalOrders(startDate, endDate)
	if err != nil || totalOrders == 0 {
		return 0, fmt.Errorf("unable to calculate AOV due to error in fetching")
	}
	averageOrderValue := totalRevenue / float64(totalOrders)
	return averageOrderValue, nil
}
func NewService(db *gorm.DB) (Service, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}
	return &Conn{db: db}, nil
}
