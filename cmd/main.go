package main

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	pageMargin = 40.0
	taxRate    = 0.10
)

type LineItem struct {
	Item        string
	Description string
	Quantity    int
	Price       int
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

type InvoicePDFGenerator struct {
	pdf    *gofpdf.Fpdf
	width  float64
	height float64
}

func NewInvoicePDFGenerator() *InvoicePDFGenerator {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	w, h := pdf.GetPageSize()
	pdf.AddPage()

	return &InvoicePDFGenerator{
		pdf:    pdf,
		width:  w,
		height: h,
	}
}

func (g *InvoicePDFGenerator) Generate(invoice *Invoice, outputPath string) error {
	g.pdf.SetFillColor(0, 0, 0)
	g.pdf.Rect(0, 0, g.width, g.height, "F")

	contentX := pageMargin
	contentY := pageMargin
	contentW := g.width - (pageMargin * 2)
	contentH := g.height - (pageMargin * 2)

	g.pdf.SetFillColor(17, 17, 17)
	g.pdf.RoundedRect(contentX, contentY, contentW, contentH, 6, "1234", "F")

	y := contentY + 30.0

	y = g.renderHeader(invoice, contentX, y, contentW)
	y = g.renderBillToSection(invoice, contentX, y, contentW)
	y = g.renderTable(invoice, contentX, y, contentW)
	_ = g.renderTotals(invoice, contentX, y, contentW)
	g.renderFooter(contentX, contentY+contentH-50.0, contentW)

	return g.pdf.OutputFileAndClose(outputPath)
}

func (g *InvoicePDFGenerator) renderHeader(invoice *Invoice, x, y, width float64) float64 {
	g.pdf.SetFont("helvetica", "B", 24)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.Text(x+20, y, "Invoice")

	g.pdf.SetFont("helvetica", "", 11)
	g.pdf.SetTextColor(115, 115, 115)
	g.pdf.Text(x+20, y+20, "Invoice #"+invoice.Number)

	rightX := x + width - 20.0
	g.pdf.SetFont("helvetica", "B", 12)
	g.pdf.SetTextColor(255, 255, 255)
	companyW := g.pdf.GetStringWidth(invoice.Company.Name)
	g.pdf.Text(rightX-companyW, y, invoice.Company.Name)

	g.pdf.SetFont("helvetica", "", 10)
	g.pdf.SetTextColor(163, 163, 163)

	addressW := g.pdf.GetStringWidth(invoice.Company.Address)
	g.pdf.Text(rightX-addressW, y+16, invoice.Company.Address)

	cityW := g.pdf.GetStringWidth(invoice.Company.City)
	g.pdf.Text(rightX-cityW, y+30, invoice.Company.City)

	phoneW := g.pdf.GetStringWidth("Phone: " + invoice.Company.Phone)
	g.pdf.Text(rightX-phoneW, y+44, "Phone: "+invoice.Company.Phone)

	emailW := g.pdf.GetStringWidth("Email: " + invoice.Company.Email)
	g.pdf.Text(rightX-emailW, y+58, "Email: "+invoice.Company.Email)

	return y + 90.0
}

func (g *InvoicePDFGenerator) renderBillToSection(invoice *Invoice, x, y, width float64) float64 {
	g.pdf.SetFont("helvetica", "B", 14)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.Text(x+20, y, "Bill To:")

	g.pdf.SetFont("helvetica", "", 10)
	g.pdf.SetTextColor(163, 163, 163)
	g.pdf.Text(x+20, y+18, invoice.Client.Name)
	g.pdf.Text(x+20, y+32, invoice.Client.Address)
	g.pdf.Text(x+20, y+46, invoice.Client.City)

	rightX := x + width - 20.0

	g.pdf.SetFont("helvetica", "B", 14)
	g.pdf.SetTextColor(255, 255, 255)
	dateLabel := "Invoice Date:"
	dateLabelW := g.pdf.GetStringWidth(dateLabel)
	g.pdf.Text(rightX-dateLabelW, y, dateLabel)

	g.pdf.SetFont("helvetica", "", 10)
	g.pdf.SetTextColor(163, 163, 163)
	dateText := invoice.Date.Format("January 2, 2006")
	dateW := g.pdf.GetStringWidth(dateText)
	g.pdf.Text(rightX-dateW, y+18, dateText)

	g.pdf.SetFont("helvetica", "B", 14)
	g.pdf.SetTextColor(255, 255, 255)
	dueLabelW := g.pdf.GetStringWidth("Due Date:")
	g.pdf.Text(rightX-dueLabelW, y+40, "Due Date:")

	g.pdf.SetFont("helvetica", "", 10)
	g.pdf.SetTextColor(163, 163, 163)
	dueText := invoice.DueDate.Format("January 2, 2006")
	dueW := g.pdf.GetStringWidth(dueText)
	g.pdf.Text(rightX-dueW, y+58, dueText)

	return y + 80.0
}

func (g *InvoicePDFGenerator) renderTable(invoice *Invoice, x, y, width float64) float64 {
	tableX := x + 20.0
	tableW := width - 40.0

	headerH := 32.0
	g.pdf.SetFillColor(38, 38, 38)
	g.pdf.Rect(tableX, y, tableW, headerH, "F")

	g.pdf.SetFont("helvetica", "B", 10)
	g.pdf.SetTextColor(163, 163, 163)

	colItem := tableX + 10
	colDesc := tableX + 80
	colQty := tableX + tableW - 210
	colPrice := tableX + tableW - 140
	colTotal := tableX + tableW - 70

	g.pdf.Text(colItem, y+20, "Item")
	g.pdf.Text(colDesc, y+20, "Description")

	qtyText := "Quantity"
	qtyW := g.pdf.GetStringWidth(qtyText)
	g.pdf.Text(colQty+70-qtyW, y+20, qtyText)

	priceText := "Price"
	priceW := g.pdf.GetStringWidth(priceText)
	g.pdf.Text(colPrice+70-priceW, y+20, priceText)

	totalText := "Total"
	totalW := g.pdf.GetStringWidth(totalText)
	g.pdf.Text(colTotal+70-totalW, y+20, totalText)

	y += headerH

	g.pdf.SetFont("helvetica", "", 10)
	rowH := 42.0

	for _, item := range invoice.LineItems {
		g.pdf.SetDrawColor(38, 38, 38)
		g.pdf.SetLineWidth(0.3)
		g.pdf.Line(tableX, y+rowH, tableX+tableW, y+rowH)

		g.pdf.SetTextColor(255, 255, 255)
		g.pdf.Text(colItem, y+18, item.Item)

		g.pdf.SetTextColor(163, 163, 163)
		g.pdf.Text(colDesc, y+18, item.Description)

		qtyStr := fmt.Sprintf("%d", item.Quantity)
		qtyStrW := g.pdf.GetStringWidth(qtyStr)
		g.pdf.Text(colQty+70-qtyStrW, y+18, qtyStr)

		priceStr := toUSD(item.Price)
		priceStrW := g.pdf.GetStringWidth(priceStr)
		g.pdf.Text(colPrice+70-priceStrW, y+18, priceStr)

		totalStr := toUSD(item.Price * item.Quantity)
		totalStrW := g.pdf.GetStringWidth(totalStr)
		g.pdf.SetTextColor(255, 255, 255)
		g.pdf.Text(colTotal+70-totalStrW, y+18, totalStr)

		y += rowH
	}

	return y + 20.0
}

func (g *InvoicePDFGenerator) renderTotals(invoice *Invoice, x, y, width float64) float64 {
	rightX := x + width - 20.0

	g.pdf.SetFont("helvetica", "B", 12)
	g.pdf.SetTextColor(255, 255, 255)

	g.pdf.Text(rightX-200, y, "Subtotal:")

	g.pdf.SetFont("helvetica", "", 12)
	g.pdf.SetTextColor(163, 163, 163)
	subtotalValue := toUSD(invoice.CalculateSubtotal())
	subtotalValueW := g.pdf.GetStringWidth(subtotalValue)
	g.pdf.Text(rightX-subtotalValueW, y, subtotalValue)

	y += 22.0

	g.pdf.SetFont("helvetica", "B", 12)
	g.pdf.SetTextColor(255, 255, 255)
	taxLabel := fmt.Sprintf("Tax (%.0f%%):", invoice.TaxRate*100)
	g.pdf.Text(rightX-200, y, taxLabel)

	g.pdf.SetFont("helvetica", "", 12)
	g.pdf.SetTextColor(163, 163, 163)
	taxValue := toUSD(invoice.CalculateTax())
	taxValueW := g.pdf.GetStringWidth(taxValue)
	g.pdf.Text(rightX-taxValueW, y, taxValue)

	y += 28.0

	g.pdf.SetFont("helvetica", "B", 12)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.Text(rightX-200, y, "Total:")

	g.pdf.SetFont("helvetica", "B", 20)
	totalValue := toUSD(invoice.CalculateTotal())
	totalValueW := g.pdf.GetStringWidth(totalValue)
	g.pdf.Text(rightX-totalValueW, y, totalValue)

	return y + 30.0
}

func (g *InvoicePDFGenerator) renderFooter(x, y, width float64) {
	g.pdf.SetFont("helvetica", "", 9)
	g.pdf.SetTextColor(115, 115, 115)

	g.pdf.Text(x+20, y, "Thank you for your business!")
	g.pdf.Text(x+20, y+16, "Please make payment by the due date to avoid late fees.")
}

func toUSD(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}
	dollars := cents / 100

	if dollars >= 1000 {
		thousands := dollars / 1000
		hundreds := dollars % 1000
		return fmt.Sprintf("$%d,%03d.%s", thousands, hundreds, centsStr)
	}

	return fmt.Sprintf("$%d.%s", dollars, centsStr)
}

func createSampleInvoice() *Invoice {
	return &Invoice{
		Number:  "123456",
		Date:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		DueDate: time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
		Company: CompanyInfo{
			Name:    "Acme Inc.",
			Address: "123 Main St.",
			City:    "Anytown, USA 12345",
			Phone:   "(123) 456-7890",
			Email:   "info@acme.com",
		},
		Client: ClientInfo{
			Name:    "John Doe",
			Address: "456 Oak St.",
			City:    "Anytown, USA 54321",
		},
		LineItems: []LineItem{
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

func main() {
	invoice := createSampleInvoice()
	generator := NewInvoicePDFGenerator()
	if err := generator.Generate(invoice, "invoice.pdf"); err != nil {
		panic(err)
	}
	fmt.Println("✅ Invoice generated successfully")
}
