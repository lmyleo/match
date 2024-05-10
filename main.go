package main

import (
	"fmt"
	"match/process"
	"os"
)

func main() {
	convNum, maxChooseNum, outputType, err := input()
	if err != nil {
		return
	}
	// convNum, maxChooseNum := 3, 3

	p := process.NewBaseProcessor(convNum, maxChooseNum)

	p.Prepare()

	p.LoadData()

	p.Match()

	if outputType == 1 {
		p.OutputData()
	} else {
		p.OutputDebug()
	}

	pause()
}

func input() (a, b, c int, err error) {
	// 读取第一个整数
	fmt.Print("请输入对话轮数（输入数字后按下回车键）: ")
	if _, err = fmt.Scanln(&a); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()

	// 读取第二个整数
	fmt.Print("请输入最大可选择人数（输入数字后按下回车键）: ")
	if _, err = fmt.Scanln(&b); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()

	// 读取第三个整数
	fmt.Println("请选择数据输出模式（输入数字后按下回车键）")
	fmt.Println("1. 输出桌号")
	fmt.Println("2. 输出匹配类型")
	if _, err = fmt.Scanln(&c); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()
	return
}

func pause() {
	fmt.Printf("\n\n按任意键退出...")

	b := make([]byte, 1)

	os.Stdin.Read(b)
}
