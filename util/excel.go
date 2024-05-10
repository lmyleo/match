package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

const (
	excelSuffix = ".xlsx"
	ouputName   = "output"
)

func ReadExcel() ([][]string, error) {
	// 打开文件
	file, err := excelize.OpenFile(getExcelFile())
	if err != nil {
		return nil, err
	}

	// 读取第一个工作表
	sheetName := file.GetSheetName(0)
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getExcelFile() string {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 读取当前目录下的所有文件和文件夹
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 找出目标 excel 文件名
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, excelSuffix) && !strings.Contains(name, ouputName) {
			return name
		}
	}
	return ""
}
