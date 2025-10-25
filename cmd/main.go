package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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

type PDFWriter struct {
	buffer      bytes.Buffer
	objects     []int
	objectCount int
	pageObjID   int
	catalogID   int
	pagesID     int
	fontObjIDs  map[string]int
}

func NewPDFWriter() *PDFWriter {
	return &PDFWriter{
		fontObjIDs: make(map[string]int),
	}
}

func (w *PDFWriter) writeString(s string) {
	w.buffer.WriteString(s)
}

func (w *PDFWriter) startObject() int {
	w.objectCount++
	w.objects = append(w.objects, w.buffer.Len())
	w.writeString(fmt.Sprintf("%d 0 obj\n", w.objectCount))
	return w.objectCount
}

func (w *PDFWriter) endObject() {
	w.writeString("endobj\n")
}

func (w *PDFWriter) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write(w.buffer.Bytes())
	return int64(n), err
}

type InvoicePDFGenerator struct {
	writer *PDFWriter
	width  float64
	height float64
}

func NewInvoicePDFGenerator() *InvoicePDFGenerator {
	// Letter size in points: 612 x 792
	return &InvoicePDFGenerator{
		writer: NewPDFWriter(),
		width:  612,
		height: 792,
	}
}

func (g *InvoicePDFGenerator) Generate(invoice *Invoice, outputPath string) error {
	// Build the PDF content stream
	content := g.buildContentStream(invoice)

	// Start building PDF
	g.writer.writeString("%PDF-1.4\n")
	g.writer.writeString("%âãÏÓ\n") // Binary comment for compatibility

	// Create font objects
	helveticaID := g.addFont("Helvetica")
	helveticaBoldID := g.addFont("Helvetica-Bold")

	// Create content stream object
	contentID := g.writer.startObject()
	g.writer.writeString(fmt.Sprintf("<< /Length %d >>\n", len(content)))
	g.writer.writeString("stream\n")
	g.writer.writeString(content)
	g.writer.writeString("\nendstream\n")
	g.writer.endObject()

	// Create page object
	g.writer.pageObjID = g.writer.startObject()
	g.writer.writeString("<< /Type /Page\n")
	g.writer.writeString("   /Parent 2 0 R\n")
	g.writer.writeString(fmt.Sprintf("   /MediaBox [0 0 %.2f %.2f]\n", g.width, g.height))
	g.writer.writeString(fmt.Sprintf("   /Contents %d 0 R\n", contentID))
	g.writer.writeString("   /Resources << /Font << ")
	g.writer.writeString(fmt.Sprintf("/F1 %d 0 R /F2 %d 0 R", helveticaID, helveticaBoldID))
	g.writer.writeString(" >> >>\n")
	g.writer.writeString(">>\n")
	g.writer.endObject()

	// Create pages object
	g.writer.pagesID = g.writer.startObject()
	g.writer.writeString("<< /Type /Pages\n")
	g.writer.writeString(fmt.Sprintf("   /Kids [%d 0 R]\n", g.writer.pageObjID))
	g.writer.writeString("   /Count 1\n")
	g.writer.writeString(">>\n")
	g.writer.endObject()

	// Create catalog object
	g.writer.catalogID = g.writer.startObject()
	g.writer.writeString("<< /Type /Catalog\n")
	g.writer.writeString(fmt.Sprintf("   /Pages %d 0 R\n", g.writer.pagesID))
	g.writer.writeString(">>\n")
	g.writer.endObject()

	// Write cross-reference table
	xrefPos := g.writer.buffer.Len()
	g.writer.writeString("xref\n")
	g.writer.writeString(fmt.Sprintf("0 %d\n", g.writer.objectCount+1))
	g.writer.writeString("0000000000 65535 f \n")
	for _, offset := range g.writer.objects {
		g.writer.writeString(fmt.Sprintf("%010d 00000 n \n", offset))
	}

	// Write trailer
	g.writer.writeString("trailer\n")
	g.writer.writeString("<<\n")
	g.writer.writeString(fmt.Sprintf("  /Size %d\n", g.writer.objectCount+1))
	g.writer.writeString(fmt.Sprintf("  /Root %d 0 R\n", g.writer.catalogID))
	g.writer.writeString(">>\n")
	g.writer.writeString("startxref\n")
	g.writer.writeString(fmt.Sprintf("%d\n", xrefPos))
	g.writer.writeString("%%EOF\n")

	// Write to file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = g.writer.WriteTo(file)
	return err
}

func (g *InvoicePDFGenerator) addFont(fontName string) int {
	objID := g.writer.startObject()
	g.writer.writeString("<< /Type /Font\n")
	g.writer.writeString("   /Subtype /Type1\n")
	g.writer.writeString(fmt.Sprintf("   /BaseFont /%s\n", fontName))
	g.writer.writeString(">>\n")
	g.writer.endObject()
	g.writer.fontObjIDs[fontName] = objID
	return objID
}

func (g *InvoicePDFGenerator) buildContentStream(invoice *Invoice) string {
	var b strings.Builder

	contentX := pageMargin
	contentY := pageMargin
	contentW := g.width - (pageMargin * 2)

	// Draw black background
	b.WriteString("0 0 0 rg\n")
	b.WriteString(fmt.Sprintf("0 0 %.2f %.2f re f\n", g.width, g.height))

	y := contentY + 30.0

	y += 90.0                                   // Header
	y += 80.0                                   // Bill To Section
	y += 32.0                                   // Table header
	y += float64(len(invoice.LineItems)) * 42.0 // Table rows
	y += 20.0                                   // Table bottom spacing
	y += 80.0                                   // Totals section
	y += 50.0                                   // Footer spacing

	contentH := y - contentY + 30.0 // Add padding at bottom

	// Draw dark gray content area with rounded corners
	g.addRoundedRect(&b, contentX, contentY, contentW, contentH, 6, 0, 0, 0)

	y = contentY + 30.0
	y = g.renderHeader(&b, invoice, contentX, y, contentW)
	y = g.renderBillToSection(&b, invoice, contentX, y, contentW)
	y += 30.0
	y = g.renderTable(&b, invoice, contentX, y, contentW)
	_ = g.renderTotals(&b, invoice, contentX, y, contentW)
	g.renderFooter(&b, contentX, g.height-50.0, contentW)

	return b.String()
}

func (g *InvoicePDFGenerator) addRoundedRect(b *strings.Builder, x, y, w, h, r float64, red, green, blue int) {
	radius := r

	// Set fill color
	fmt.Fprintf(b, "%.3f %.3f %.3f rg\n", float64(red)/255.0, float64(green)/255.0, float64(blue)/255.0)

	// Calculate control point offset for circular arc approximation
	// Using magic number 0.552284749831 for quarter circle
	cp := radius * 0.552284749831

	x1 := x
	y1 := y
	x2 := x + w
	y2 := y + h

	// Start from bottom-left corner (after radius)
	fmt.Fprintf(b, "%.2f %.2f m\n", x1+radius, y1)

	// Bottom edge
	fmt.Fprintf(b, "%.2f %.2f l\n", x2-radius, y1)

	// Bottom-right corner
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n",
		x2-radius+cp, y1, x2, y1+radius-cp, x2, y1+radius)

	// Right edge
	fmt.Fprintf(b, "%.2f %.2f l\n", x2, y2-radius)

	// Top-right corner
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n",
		x2, y2-radius+cp, x2-radius+cp, y2, x2-radius, y2)

	// Top edge
	fmt.Fprintf(b, "%.2f %.2f l\n", x1+radius, y2)

	// Top-left corner
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n",
		x1+radius-cp, y2, x1, y2-radius+cp, x1, y2-radius)

	// Left edge
	fmt.Fprintf(b, "%.2f %.2f l\n", x1, y1+radius)

	// Bottom-left corner
	fmt.Fprintf(b, "%.2f %.2f %.2f %.2f %.2f %.2f c\n",
		x1, y1+radius-cp, x1+radius-cp, y1, x1+radius, y1)

	// Fill the path
	fmt.Fprintf(b, "f\n")
}

func (g *InvoicePDFGenerator) renderHeader(b *strings.Builder, invoice *Invoice, x, y, width float64) float64 {
	g.setText(b, "F2", 24, 255, 255, 255)
	g.drawText(b, x+20, y, "Invoice")

	g.setText(b, "F1", 11, 115, 115, 115)
	g.drawText(b, x+20, y+15, "Invoice #"+invoice.Number)

	rightX := x + width - 20.0

	g.setText(b, "F2", 12, 255, 255, 255)
	companyW := g.getStringWidth(invoice.Company.Name, 12)
	g.drawText(b, rightX-companyW, y-7, invoice.Company.Name)

	g.setText(b, "F1", 10, 163, 163, 163)

	addressW := g.getStringWidth(invoice.Company.Address, 10)
	g.drawText(b, rightX-addressW, y+10, invoice.Company.Address)

	cityW := g.getStringWidth(invoice.Company.City, 10)
	g.drawText(b, rightX-cityW, y+23, invoice.Company.City)

	phoneText := "Phone: " + invoice.Company.Phone
	phoneW := g.getStringWidth(phoneText, 10)
	g.drawText(b, rightX-phoneW, y+37, phoneText)

	emailText := "Email: " + invoice.Company.Email
	emailW := g.getStringWidth(emailText, 10)
	g.drawText(b, rightX-emailW, y+51, emailText)

	return y + 90.0
}

func (g *InvoicePDFGenerator) renderBillToSection(b *strings.Builder, invoice *Invoice, x, y, width float64) float64 {
	g.setText(b, "F2", 14, 255, 255, 255)
	g.drawText(b, x+20, y, "Bill To:")

	g.setText(b, "F1", 10, 163, 163, 163)
	g.drawText(b, x+20, y+18, invoice.Client.Name)
	g.drawText(b, x+20, y+32, invoice.Client.Address)
	g.drawText(b, x+20, y+46, invoice.Client.City)

	rightX := x + width - 20.0

	g.setText(b, "F2", 14, 255, 255, 255)
	dateLabel := "Invoice Date:"
	dateLabelW := g.getStringWidth(dateLabel, 14)
	g.drawText(b, rightX-dateLabelW, y, dateLabel)

	g.setText(b, "F1", 10, 163, 163, 163)
	dateText := invoice.Date.Format("January 2, 2006")
	dateW := g.getStringWidth(dateText, 10)
	g.drawText(b, rightX-dateW, y+18, dateText)

	g.setText(b, "F2", 14, 255, 255, 255)
	dueLabel := "Due Date:"
	dueLabelW := g.getStringWidth(dueLabel, 14)
	g.drawText(b, rightX-dueLabelW, y+40, dueLabel)

	g.setText(b, "F1", 10, 163, 163, 163)
	dueText := invoice.DueDate.Format("January 2, 2006")
	dueW := g.getStringWidth(dueText, 10)
	g.drawText(b, rightX-dueW, y+58, dueText)

	return y + 80.0
}

func (g *InvoicePDFGenerator) renderTable(b *strings.Builder, invoice *Invoice, x, y, width float64) float64 {
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
	qtyW := g.getStringWidth(qtyText, 10)
	g.drawText(b, colQty+70-qtyW, y+20, qtyText)

	priceText := "Price"
	priceW := g.getStringWidth(priceText, 10)
	g.drawText(b, colPrice+70-priceW, y+20, priceText)

	totalText := "Total"
	totalW := g.getStringWidth(totalText, 10)
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
		qtyStrW := g.getStringWidth(qtyStr, 10)
		g.drawText(b, colQty+70-qtyStrW, y+18, qtyStr)

		// Price (right-aligned)
		priceStr := toUSD(item.Price)
		priceStrW := g.getStringWidth(priceStr, 10)
		g.drawText(b, colPrice+70-priceStrW, y+18, priceStr)

		// Total (right-aligned, white)
		g.setText(b, "F1", 10, 255, 255, 255)
		totalStr := toUSD(item.Price * item.Quantity)
		totalStrW := g.getStringWidth(totalStr, 10)
		g.drawText(b, colTotal+70-totalStrW, y+18, totalStr)

		y += rowH
	}

	return y + 20.0
}

func (g *InvoicePDFGenerator) renderTotals(b *strings.Builder, invoice *Invoice, x, y, width float64) float64 {
	rightX := x + width - 20.0

	g.setText(b, "F2", 10, 255, 255, 255)
	g.drawText(b, rightX-200, y, "Subtotal:")

	g.setText(b, "F1", 12, 163, 163, 163)
	subtotalValue := toUSD(invoice.CalculateSubtotal())
	subtotalValueW := g.getStringWidth(subtotalValue, 12)
	g.drawText(b, rightX-subtotalValueW, y, subtotalValue)

	y += 22.0

	g.setText(b, "F2", 10, 255, 255, 255)
	taxLabel := fmt.Sprintf("Tax (%.0f%%):", invoice.TaxRate*100)
	g.drawText(b, rightX-200, y, taxLabel)

	g.setText(b, "F1", 12, 163, 163, 163)
	taxValue := toUSD(invoice.CalculateTax())
	taxValueW := g.getStringWidth(taxValue, 12)
	g.drawText(b, rightX-taxValueW, y, taxValue)

	y += 48.0

	g.setText(b, "F2", 20, 255, 255, 255)
	g.drawText(b, rightX-200, y, "Total:")

	g.setText(b, "F2", 20, 255, 255, 255)
	totalValue := toUSD(invoice.CalculateTotal())
	totalValueW := g.getStringWidth(totalValue, 20)
	g.drawText(b, rightX-totalValueW, y, totalValue)

	return y + 30.0
}

func (g *InvoicePDFGenerator) renderFooter(b *strings.Builder, x, y, width float64) {
	g.setText(b, "F1", 9, 115, 115, 115)
	g.drawText(b, x+20, y, "Thank you for your business!")
	g.drawText(b, x+20, y+16, "Please make payment by the due date to avoid late fees.")
}

func (g *InvoicePDFGenerator) setText(b *strings.Builder, font string, size int, r, gr, bl int) {
	fmt.Fprintf(b, "/%s %d Tf\n", font, size)
	fmt.Fprintf(b, "%.3f %.3f %.3f rg\n", float64(r)/255.0, float64(gr)/255.0, float64(bl)/255.0)
}

func (g *InvoicePDFGenerator) drawText(b *strings.Builder, x, y float64, text string) {
	yPDF := g.height - y

	b.WriteString("BT\n")
	fmt.Fprintf(b, "%.2f %.2f Td\n", x, yPDF)
	fmt.Fprintf(b, "(%s) Tj\n", escapePDFString(text))
	b.WriteString("ET\n")
}

// escapePDFString escapes special characters in PDF strings
func escapePDFString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}

// getStringWidth estimates the width of a string in points
// Based on standard Helvetica character widths
func (g *InvoicePDFGenerator) getStringWidth(s string, fontSize int) float64 {
	// Standard Helvetica character widths (relative to 1000 units)
	widths := map[rune]int{
		' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889, '&': 667, '\'': 191,
		'(': 333, ')': 333, '*': 389, '+': 584, ',': 278, '-': 333, '.': 278, '/': 278,
		'0': 556, '1': 556, '2': 556, '3': 556, '4': 556, '5': 556, '6': 556, '7': 556,
		'8': 556, '9': 556, ':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556,
		'@': 1015, 'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778,
		'H': 722, 'I': 278, 'J': 500, 'K': 667, 'L': 556, 'M': 833, 'N': 722, 'O': 778,
		'P': 667, 'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944,
		'X': 667, 'Y': 667, 'Z': 611, '[': 278, '\\': 278, ']': 278, '^': 469, '_': 556,
		'`': 333, 'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278, 'g': 556,
		'h': 556, 'i': 222, 'j': 222, 'k': 500, 'l': 222, 'm': 833, 'n': 556, 'o': 556,
		'p': 556, 'q': 556, 'r': 333, 's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722,
		'x': 500, 'y': 500, 'z': 500, '{': 334, '|': 260, '}': 334, '~': 584,
	}

	totalWidth := 0
	for _, ch := range s {
		if w, ok := widths[ch]; ok {
			totalWidth += w
		} else {
			totalWidth += 556
		}
	}

	// Convert from 1000-unit scale to actual points
	return float64(totalWidth) * float64(fontSize) / 1000.0
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
