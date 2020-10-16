package imagestopdf

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func ConvertToPdf(pdfFileName string, images []string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, imgPath := range images {
		pdf.AddPage()
		imgOptions := gofpdf.ImageOptions{
			AllowNegativePosition: true,
			ImageType:             strings.Split(imgPath, ".")[1],
			ReadDpi:               true,
		}
		pdf.RegisterImageOptions(imgPath, imgOptions)
		pdf.ImageOptions(imgPath, 0, 0, 200, 290, false, imgOptions, 0, "")
	}
	return pdf.OutputFileAndClose(pdfFileName)
}

func CompressPdf(pdfFileName string) (string, error) {
	pdfConfiguration := pdfcpu.NewDefaultConfiguration()
	outputFileName := fmt.Sprintf("%s%s", strings.Split(pdfFileName, ".")[0], "_optimized.pdf")
	return outputFileName, pdfcpuapi.OptimizeFile(pdfFileName, outputFileName, pdfConfiguration)
}
