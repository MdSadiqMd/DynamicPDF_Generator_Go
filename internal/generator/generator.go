package generator

import (
	"os"
	"strings"

	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/models"
	"github.com/MdSadiqMd/DynamicPDF_Generator_Go/internal/pdf"
)

const (
	PageMargin = 40.0
)

type InvoiceGenerator struct {
	writer *pdf.Writer
	width  float64
	height float64
}

func New() *InvoiceGenerator {
	// Letter size in points: 612 x 792
	return &InvoiceGenerator{
		writer: pdf.NewWriter(),
		width:  612,
		height: 792,
	}
}

func (g *InvoiceGenerator) Generate(invoice *models.Invoice, outputPath string) error {
	// Build the PDF content stream
	content := g.buildContentStream(invoice)

	// Build the complete PDF with fonts
	fonts := []string{"Helvetica", "Helvetica-Bold"}
	if err := g.writer.Build(content, fonts, g.width, g.height); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = g.writer.WriteTo(file)
	return err
}

func (g *InvoiceGenerator) buildContentStream(invoice *models.Invoice) string {
	var b strings.Builder

	contentX := PageMargin
	contentY := PageMargin
	contentW := g.width - (PageMargin * 2)

	// Draw black background
	b.WriteString("0 0 0 rg\n")
	b.WriteString("0 0 612.00 792.00 re f\n")

	// Calculate content height
	y := contentY + 30.0
	y += 90.0                                   // Header
	y += 80.0                                   // Bill To Section
	y += 32.0                                   // Table header
	y += float64(len(invoice.LineItems)) * 42.0 // Table rows
	y += 20.0                                   // Table bottom spacing
	y += 80.0                                   // Totals section
	y += 50.0                                   // Footer spacing
	contentH := y - contentY + 30.0

	// Draw dark gray content area with rounded corners
	g.addRoundedRect(&b, contentX, contentY, contentW, contentH, 6, 0, 0, 0)

	// Render sections
	y = contentY + 30.0
	y = g.renderHeader(&b, invoice, contentX, y, contentW)
	y = g.renderBillToSection(&b, invoice, contentX, y, contentW)
	y += 30.0
	y = g.renderTable(&b, invoice, contentX, y, contentW)
	_ = g.renderTotals(&b, invoice, contentX, y, contentW)
	g.renderFooter(&b, contentX, g.height-50.0, contentW)

	return b.String()
}
