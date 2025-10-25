package main

import (
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	w, h := pdf.GetPageSize()
	fmt.Println(w, h)

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 38)
	_, lineHeight := pdf.GetFontSize()
	pdf.Cell(w, lineHeight, "Hello World!")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("PDF generated successfully")
	}
}
