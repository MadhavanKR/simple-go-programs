package main

import (
	"fmt"
	"os"
	"strings"

	"./imagestopdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("insufficient arguments. format = go run imagesToPdfMain.go action")
		os.Exit(1)
	}

	action := os.Args[1]
	if strings.EqualFold(action, "convert") {
		if len(os.Args) < 4 {
			fmt.Println("insufficient arguments. format = go run imagesToPdfMain.go convert pdfFilePath img1 img2 ..")
			os.Exit(1)
		}
		pdfFilePath := os.Args[2]
		imageFiles := os.Args[3:len(os.Args)]
		convertImgToPdfError := imagestopdf.ConvertToPdf(pdfFilePath, imageFiles)
		if convertImgToPdfError != nil {
			fmt.Println("failed to create pdf from images: ", convertImgToPdfError)
			os.Exit(1)
		}
		fmt.Println("successfully created: ", pdfFilePath)
	} else if strings.EqualFold(action, "optimize") {
		if len(os.Args) < 3 {
			fmt.Println("insufficient arguments. format = go run imagesToPdfMain.go optimize pdfFilePath")
			os.Exit(1)
		}
		pdfFilePath := os.Args[2]
		outputPdfFilePath, optimizePdfError := imagestopdf.CompressPdf(pdfFilePath)
		if optimizePdfError != nil {
			fmt.Println("failed to create pdf from images: ", optimizePdfError)
			os.Exit(1)
		}
		fmt.Println("successfully created: ", outputPdfFilePath)
	}

}
