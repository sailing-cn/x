package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func test() {
	_, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
}
