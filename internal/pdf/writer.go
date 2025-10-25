package pdf

import (
	"bytes"
	"fmt"
	"io"
)

// Writer handles low-level PDF file format operations
type Writer struct {
	buffer      bytes.Buffer
	objects     []int // byte offsets of each object
	objectCount int
	pageObjID   int
	catalogID   int
	pagesID     int
	fontObjIDs  map[string]int // font name to object ID
}

// NewWriter creates a new PDF writer
func NewWriter() *Writer {
	return &Writer{
		fontObjIDs: make(map[string]int),
	}
}

// WriteString writes a string to the PDF buffer
func (w *Writer) WriteString(s string) {
	w.buffer.WriteString(s)
}

// StartObject begins a new PDF object and returns its ID
func (w *Writer) StartObject() int {
	w.objectCount++
	w.objects = append(w.objects, w.buffer.Len())
	w.WriteString(fmt.Sprintf("%d 0 obj\n", w.objectCount))
	return w.objectCount
}

// EndObject closes the current PDF object
func (w *Writer) EndObject() {
	w.WriteString("endobj\n")
}

// WriteTo writes the PDF buffer to the given writer
func (w *Writer) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write(w.buffer.Bytes())
	return int64(n), err
}

// Build constructs the complete PDF file with header, objects, xref, and trailer
func (w *Writer) Build(content string, fonts []string, pageWidth, pageHeight float64) error {
	// PDF header
	w.WriteString("%PDF-1.4\n")
	w.WriteString("%âãÏÓ\n") // Binary comment for compatibility

	// Create font objects
	fontRefs := make(map[string]int)
	for _, fontName := range fonts {
		objID := w.StartObject()
		w.WriteString("<< /Type /Font\n")
		w.WriteString("   /Subtype /Type1\n")
		w.WriteString(fmt.Sprintf("   /BaseFont /%s\n", fontName))
		w.WriteString(">>\n")
		w.EndObject()
		fontRefs[fontName] = objID
	}

	// Create content stream object
	contentID := w.StartObject()
	w.WriteString(fmt.Sprintf("<< /Length %d >>\n", len(content)))
	w.WriteString("stream\n")
	w.WriteString(content)
	w.WriteString("\nendstream\n")
	w.EndObject()

	// Create page object
	w.pageObjID = w.StartObject()
	w.WriteString("<< /Type /Page\n")
	w.WriteString("   /Parent 2 0 R\n")
	w.WriteString(fmt.Sprintf("   /MediaBox [0 0 %.2f %.2f]\n", pageWidth, pageHeight))
	w.WriteString(fmt.Sprintf("   /Contents %d 0 R\n", contentID))
	w.WriteString("   /Resources << /Font << ")

	// Add font references
	for i, fontName := range fonts {
		w.WriteString(fmt.Sprintf("/F%d %d 0 R", i+1, fontRefs[fontName]))
		if i < len(fonts)-1 {
			w.WriteString(" ")
		}
	}
	w.WriteString(" >> >>\n")
	w.WriteString(">>\n")
	w.EndObject()

	// Create pages object
	w.pagesID = w.StartObject()
	w.WriteString("<< /Type /Pages\n")
	w.WriteString(fmt.Sprintf("   /Kids [%d 0 R]\n", w.pageObjID))
	w.WriteString("   /Count 1\n")
	w.WriteString(">>\n")
	w.EndObject()

	// Create catalog object
	w.catalogID = w.StartObject()
	w.WriteString("<< /Type /Catalog\n")
	w.WriteString(fmt.Sprintf("   /Pages %d 0 R\n", w.pagesID))
	w.WriteString(">>\n")
	w.EndObject()

	// Write cross-reference table
	xrefPos := w.buffer.Len()
	w.WriteString("xref\n")
	w.WriteString(fmt.Sprintf("0 %d\n", w.objectCount+1))
	w.WriteString("0000000000 65535 f \n")
	for _, offset := range w.objects {
		w.WriteString(fmt.Sprintf("%010d 00000 n \n", offset))
	}

	// Write trailer
	w.WriteString("trailer\n")
	w.WriteString("<<\n")
	w.WriteString(fmt.Sprintf("  /Size %d\n", w.objectCount+1))
	w.WriteString(fmt.Sprintf("  /Root %d 0 R\n", w.catalogID))
	w.WriteString(">>\n")
	w.WriteString("startxref\n")
	w.WriteString(fmt.Sprintf("%d\n", xrefPos))
	w.WriteString("%%EOF\n")

	return nil
}
