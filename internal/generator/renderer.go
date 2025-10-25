package generator

import (
	"fmt"
	"strings"

	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/models"
	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/utils"
)

// addRoundedRect creates a rounded rectangle using Bezier curves
func (g *InvoiceGenerator) addRoundedRect(b *strings.Builder, x, y, w, h, r float64, red, green, blue int) {
	radius := r

	// Set fill color
	fmt.Fprintf(b, "%.3f %.3f %.3f rg\n", float64(red)/255.0, float64(green)/255.0, float64(blue)/255.0)

	// Calculate control point offset for circular arc approximation
	cp := radius * 0.552284749831

	x1 := x
	y1 := y
	x2 := x + w
	y2 := y + h

	// Start from bottom-left corner (after radius)
	fmt.Fprintf(b, "%.2f %.2f m\n", x1+radius, y1)
	fmt.Fprintf(b, "%.2f %.2f l\n", x2-radius, y1)
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", x2-radius+cp, y1, x2, y1+radius-cp, x2, y1+radius)
	fmt.Fprintf(b, "%.2f %.2f l\n", x2, y2-radius)
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", x2, y2-radius+cp, x2-radius+cp, y2, x2-radius, y2)
	fmt.Fprintf(b, "%.2f %.2f l\n", x1+radius, y2)
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", x1+radius-cp, y2, x1, y2-radius+cp, x1, y2-radius)
	fmt.Fprintf(b, "%.2f %.2f l\n", x1, y1+radius)
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", x1, y1+radius-cp, x1+radius-cp, y1, x1+radius, y1)
	fmt.Fprintf(b, "f\n")
}

func (g *InvoiceGenerator) renderHeader(b *strings.Builder, invoice *models.Invoice, x, y, width float64) float64 {
	g.setText(b, "F2", 24, 255, 255, 255)
	g.drawText(b, x+20, y, "Invoice")

	g.setText(b, "F1", 11, 115, 115, 115)
	g.drawText(b, x+20, y+15, "Invoice #"+invoice.Number)

	rightX := x + width - 20.0

	g.setText(b, "F2", 12, 255, 255, 255)
	companyW := utils.GetStringWidth(invoice.Company.Name, 12)
	g.drawText(b, rightX-companyW, y-7, invoice.Company.Name)

	g.setText(b, "F1", 10, 163, 163, 163)

	addressW := utils.GetStringWidth(invoice.Company.Address, 10)
	g.drawText(b, rightX-addressW, y+10, invoice.Company.Address)

	cityW := utils.GetStringWidth(invoice.Company.City, 10)
	g.drawText(b, rightX-cityW, y+23, invoice.Company.City)

	phoneText := "Phone: " + invoice.Company.Phone
	phoneW := utils.GetStringWidth(phoneText, 10)
	g.drawText(b, rightX-phoneW, y+37, phoneText)

	emailText := "Email: " + invoice.Company.Email
	emailW := utils.GetStringWidth(emailText, 10)
	g.drawText(b, rightX-emailW, y+51, emailText)

	return y + 90.0
}

func (g *InvoiceGenerator) renderBillToSection(b *strings.Builder, invoice *models.Invoice, x, y, width float64) float64 {
	g.setText(b, "F2", 14, 255, 255, 255)
	g.drawText(b, x+20, y, "Bill To:")

	g.setText(b, "F1", 10, 163, 163, 163)
	g.drawText(b, x+20, y+18, invoice.Client.Name)
	g.drawText(b, x+20, y+32, invoice.Client.Address)
	g.drawText(b, x+20, y+46, invoice.Client.City)

	rightX := x + width - 20.0

	g.setText(b, "F2", 14, 255, 255, 255)
	dateLabel := "Invoice Date:"
	dateLabelW := utils.GetStringWidth(dateLabel, 14)
	g.drawText(b, rightX-dateLabelW, y, dateLabel)

	g.setText(b, "F1", 10, 163, 163, 163)
	dateText := invoice.Date.Format("January 2, 2006")
	dateW := utils.GetStringWidth(dateText, 10)
	g.drawText(b, rightX-dateW, y+18, dateText)

	g.setText(b, "F2", 14, 255, 255, 255)
	dueLabel := "Due Date:"
	dueLabelW := utils.GetStringWidth(dueLabel, 14)
	g.drawText(b, rightX-dueLabelW, y+40, dueLabel)

	g.setText(b, "F1", 10, 163, 163, 163)
	dueText := invoice.DueDate.Format("January 2, 2006")
	dueW := utils.GetStringWidth(dueText, 10)
	g.drawText(b, rightX-dueW, y+58, dueText)

	return y + 80.0
}

func (g *InvoiceGenerator) renderTable(b *strings.Builder, invoice *models.Invoice, x, y, width float64) float64 {
	tableX := x + 20.0
	tableW := width - 40.0

	headerH := 32.0

	g.setText(b, "F2", 10, 163, 163, 163)

	colItem := tableX + 10
	colDesc := tableX + 80
	colQty := tableX + tableW - 210
	colPrice := tableX + tableW - 140
	colTotal := tableX + tableW - 70

	g.drawText(b, colItem, y+20, "Item")
	g.drawText(b, colDesc, y+20, "Description")

	qtyText := "Quantity"
	qtyW := utils.GetStringWidth(qtyText, 10)
	g.drawText(b, colQty+70-qtyW, y+20, qtyText)

	priceText := "Price"
	priceW := utils.GetStringWidth(priceText, 10)
	g.drawText(b, colPrice+70-priceW, y+20, priceText)

	totalText := "Total"
	totalW := utils.GetStringWidth(totalText, 10)
	g.drawText(b, colTotal+70-totalW, y+20, totalText)

	y += headerH

	g.setText(b, "F1", 10, 163, 163, 163)
	rowH := 42.0

	for _, item := range invoice.LineItems {
		// Item name (white)
		g.setText(b, "F1", 10, 255, 255, 255)
		g.drawText(b, colItem, y+18, item.Item)

		// Description (gray)
		g.setText(b, "F1", 10, 163, 163, 163)
		g.drawText(b, colDesc, y+18, item.Description)

		// Quantity (right-aligned)
		qtyStr := fmt.Sprintf("%d", item.Quantity)
		qtyStrW := utils.GetStringWidth(qtyStr, 10)
		g.drawText(b, colQty+70-qtyStrW, y+18, qtyStr)

		// Price (right-aligned)
		priceStr := utils.ToUSD(item.Price)
		priceStrW := utils.GetStringWidth(priceStr, 10)
		g.drawText(b, colPrice+70-priceStrW, y+18, priceStr)

		// Total (right-aligned, white)
		g.setText(b, "F1", 10, 255, 255, 255)
		totalStr := utils.ToUSD(item.Price * item.Quantity)
		totalStrW := utils.GetStringWidth(totalStr, 10)
		g.drawText(b, colTotal+70-totalStrW, y+18, totalStr)

		y += rowH
	}

	return y + 20.0
}

func (g *InvoiceGenerator) renderTotals(b *strings.Builder, invoice *models.Invoice, x, y, width float64) float64 {
	rightX := x + width - 20.0

	g.setText(b, "F2", 10, 255, 255, 255)
	g.drawText(b, rightX-200, y, "Subtotal:")

	g.setText(b, "F1", 12, 163, 163, 163)
	subtotalValue := utils.ToUSD(invoice.CalculateSubtotal())
	subtotalValueW := utils.GetStringWidth(subtotalValue, 12)
	g.drawText(b, rightX-subtotalValueW, y, subtotalValue)

	y += 22.0

	g.setText(b, "F2", 10, 255, 255, 255)
	taxLabel := fmt.Sprintf("Tax (%.0f%%):", invoice.TaxRate*100)
	g.drawText(b, rightX-200, y, taxLabel)

	g.setText(b, "F1", 12, 163, 163, 163)
	taxValue := utils.ToUSD(invoice.CalculateTax())
	taxValueW := utils.GetStringWidth(taxValue, 12)
	g.drawText(b, rightX-taxValueW, y, taxValue)

	y += 48.0

	g.setText(b, "F2", 20, 255, 255, 255)
	g.drawText(b, rightX-200, y, "Total:")

	g.setText(b, "F2", 20, 255, 255, 255)
	totalValue := utils.ToUSD(invoice.CalculateTotal())
	totalValueW := utils.GetStringWidth(totalValue, 20)
	g.drawText(b, rightX-totalValueW, y, totalValue)

	return y + 30.0
}

func (g *InvoiceGenerator) renderFooter(b *strings.Builder, x, y, width float64) {
	g.setText(b, "F1", 9, 115, 115, 115)
	g.drawText(b, x+20, y, "Thank you for your business!")
	g.drawText(b, x+20, y+16, "Please make payment by the due date to avoid late fees.")
}

func (g *InvoiceGenerator) setText(b *strings.Builder, font string, size int, r, gr, bl int) {
	fmt.Fprintf(b, "/%s %d Tf\n", font, size)
	fmt.Fprintf(b, "%.3f %.3f %.3f rg\n", float64(r)/255.0, float64(gr)/255.0, float64(bl)/255.0)
}

func (g *InvoiceGenerator) drawText(b *strings.Builder, x, y float64, text string) {
	yPDF := g.height - y

	b.WriteString("BT\n")
	fmt.Fprintf(b, "%.2f %.2f Td\n", x, yPDF)
	fmt.Fprintf(b, "(%s) Tj\n", utils.EscapePDFString(text))
	b.WriteString("ET\n")
}
