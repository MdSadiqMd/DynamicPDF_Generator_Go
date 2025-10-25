# DynamicPDF_Generator_Go
A Go application that generates PDF documents entirely from scratch using only the standard library, manually implementing the PDF 1.4 specification. It defines structured data, constructs raw PDF objects, fonts, and content streams, and builds the layout by writing PDF graphics operators directly into the content stream. The core idea is to bypass external libraries like [gofpdf](https://github.com/jung-kurt/gofpdf) by manually writing PDF syntax (like rg for colors, Tf for fonts, Tj for text, re f for rectangles) directly to a byte buffer, and at last assembling all objects and references into a valid, well-formatted PDF file with accurate text alignment

## Output Example
<img width="1570" height="2042" alt="image" src="https://github.com/user-attachments/assets/b92c3bb4-46b7-4da4-9197-c16ff6ad59dd" />
