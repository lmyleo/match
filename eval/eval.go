package eval

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func genData(ids []int64, maxChoose int) map[int64]string {
	rand.New(rand.NewSource(time.Now().Unix()))
	choose := make(map[int64][]int64)
	chooseStr := make(map[int64]string)

	for _, id := range ids {
		choose[id] = make([]int64, 0)
		chooseStr[id] = ""
		chooseNum := rand.Intn(maxChoose + 1)
		for i := 0; i < chooseNum; i++ {
			v := rand.Intn(len(ids))
			if ids[v] != id {
				choose[id] = append(choose[id], ids[v])
				chooseStr[id] += fmt.Sprintf("%d，", ids[v])
			}
		}
	}

	return chooseStr
}

func genExcel(data map[int64]string, dirPath string) {
	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02_15-04-05")

	// 创建Excel文件
	f := excelize.NewFile()

	// 设置表头（可选，因为你说第一行不展示数据）
	f.SetCellValue("Sheet1", "B1", "Key")
	f.SetCellValue("Sheet1", "C1", "Value")

	// 写入数据到Excel文件
	rowIndex := 2 // 从第二行开始写入数据
	for key, value := range data {
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowIndex), key)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowIndex), value)
		rowIndex++
	}

	// 保存Excel文件到指定文件夹
	fileName := filepath.Join(dirPath, fmt.Sprintf("data_%s.xlsx", currentTime))
	err := f.SaveAs(fileName)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
}

func Evaluation(size int64, maxChoose int) {
	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02_15-04-05")

	// 创建包含格式化日期的文件夹
	dirName := fmt.Sprintf("test_%s", currentTime)
	dirPath := filepath.Join(".\\eval_test", dirName)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	ids := make([]int64, size)
	for i, _ := range ids {
		ids[i] = int64(i + 1)
	}

	data := genData(ids, maxChoose)
	genExcel(data, dirPath)
}
