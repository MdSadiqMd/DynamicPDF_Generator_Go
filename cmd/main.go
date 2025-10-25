package main

import (
	"fmt"
	"time"

	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/generator"
	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/models"
)

const (
	taxRate = 0.10
)

func main() {
	invoice := createSampleInvoice()
	gen := generator.New()
	if err := gen.Generate(invoice, "invoice.pdf"); err != nil {
		panic(err)
	}
	fmt.Println("✅ Invoice generated successfully")
}

func createSampleInvoice() *models.Invoice {
	return &models.Invoice{
		Number:  "123456",
		Date:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		DueDate: time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
		Company: models.CompanyInfo{
			Name:    "Acme Inc.",
			Address: "123 Main St.",
			City:    "Anytown, USA 12345",
			Phone:   "(123) 456-7890",
			Email:   "info@acme.com",
		},
		Client: models.ClientInfo{
			Name:    "John Doe",
			Address: "456 Oak St.",
			City:    "Anytown, USA 54321",
		},
		LineItems: []models.LineItem{
			{
				Item:        "Web Design",
				Description: "Design and development of a new website",
				Quantity:    1,
				Price:       250000, // $2,500.00
			},
			{
				Item:        "Branding",
				Description: "Logo design and brand identity",
				Quantity:    1,
				Price:       100000, // $1,000.00
			},
			{
				Item:        "Consulting",
				Description: "Strategic planning and consulting services",
				Quantity:    10,
				Price:       15000, // $150.00
			},
		},
		TaxRate: taxRate,
	}
}
