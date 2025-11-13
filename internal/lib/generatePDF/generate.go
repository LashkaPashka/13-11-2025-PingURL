package generatepdf

import (
	"bytes"
	"fmt"

	"github.com/LashkaPashka/LinkCheck/internal/model"
	"github.com/jung-kurt/gofpdf"
)

func Generate(urls []model.Url) (buf []byte, err error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "", 14)

	lines := []string{
		"Link Check Report",
		"=======================",
	}

	for _, url := range urls {
		lines = append(lines, fmt.Sprintf("https://%s - %s", url.UrlName, url.Available))
	}

	for _, line := range lines {
		pdf.CellFormat(0, 10, line, "", 1, "L", false, 0, "")
	}

	var buff bytes.Buffer
    
	err = pdf.Output(&buff)
    if err != nil {
        return buf, err
    }

	return buff.Bytes(), err
}
