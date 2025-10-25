package models

import "time"

type LineItem struct {
	Item        string
	Description string
	Quantity    int
	Price       int // In cents
}

type CompanyInfo struct {
	Name    string
	Address string
	City    string
	Phone   string
	Email   string
}

type ClientInfo struct {
	Name    string
	Address string
	City    string
}

type Invoice struct {
	Number    string
	Date      time.Time
	DueDate   time.Time
	Company   CompanyInfo
	Client    ClientInfo
	LineItems []LineItem
	TaxRate   float64
}

func (inv *Invoice) CalculateSubtotal() int {
	subtotal := 0
	for _, item := range inv.LineItems {
		subtotal += item.Price * item.Quantity
	}
	return subtotal
}

func (inv *Invoice) CalculateTax() int {
	return int(float64(inv.CalculateSubtotal()) * inv.TaxRate)
}

func (inv *Invoice) CalculateTotal() int {
	return inv.CalculateSubtotal() + inv.CalculateTax()
}
