package main

import (
	"fmt"

	"github.com/ledongthuc/pdf"
)

func pdfReader() {
	pdf.DebugOn = true

	f, r, err := pdf.Open("/home/toor/Desktop/all/Qdrant/day-11/part2/_OceanofPDF.com_The_Ultimate_Ubuntu_Handbook_-_Ken_VanDine.pdf")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Println("Total pages:", r.NumPage())

	for pageNumber := 1; pageNumber <= r.NumPage(); pageNumber++ {
		page := r.Page(pageNumber)

		if page.V.IsNull() {
			fmt.Printf("page %d is empty.", pageNumber)
			continue
		}

		content, err := page.GetPlainText(nil)

		if err != nil {
			fmt.Printf("error reading page %d:%v\n", pageNumber, err)
			continue
		}
		fmt.Printf("\n========== PDF PAGE %d ==========\n", pageNumber)
		fmt.Println(content)
	}
}

func main() {
	pdfReader()
}
