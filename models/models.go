package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	OrderID        string    `gorm:"unique"`
	CustomerID     string
	ProductID      string
	QuantitySold   int
	UnitPrice      float64
	Discount       float64
	ShippingCost   float64
	PaymentMethod  string
	DateOfSale     time.Time `gorm:"type:date"`
	ProductName    string
	Category       string
	Region         string
	CustomerName   string
	CustomerEmail  string
	CustomerAddress string
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
	Customer  Customer  `gorm:"foreignKey:CustomerID;references:ID"`
}

type Product struct {
	gorm.Model
	ID          string  `gorm:"primaryKey"`
	ProductName string
	Category    string
	Price       float64
}

type Customer struct {
	gorm.Model
	ID              string `gorm:"primaryKey"`
	Name            string
	Email           string
	Address         string
	Orders          []Order `gorm:"foreignKey:CustomerID"`
}
