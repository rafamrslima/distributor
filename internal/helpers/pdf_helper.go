package helpers

import (
	"bytes"

	"github.com/jung-kurt/gofpdf"
)

func CreatePdfReport(content string) (bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, content)

	var buf bytes.Buffer
	err := pdf.Output(&buf) // write PDF into buffer
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}
