package api

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// generateQRCodeSVG generates an SVG QR code for the given content
func generateQRCodeSVG(content string) string {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return ""
	}

	// Get the QR code bitmap
	bitmap := qr.Bitmap()
	size := len(bitmap)

	// Generate SVG
	var svg strings.Builder
	moduleSize := 8
	totalSize := size * moduleSize

	svg.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="512" height="512">`, totalSize, totalSize))
	svg.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, totalSize, totalSize))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bitmap[y][x] {
				svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="black"/>`,
					x*moduleSize, y*moduleSize, moduleSize, moduleSize))
			}
		}
	}

	svg.WriteString("</svg>")
	return svg.String()
}

