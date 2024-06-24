package main

import (
	"fmt"
	"match/eval"
	"match/process"
	"os"
	"sync"
)

const (
	randMatchRetry = 100 // 随机匹配最大重试 100 次
)

func main() {
	// algorithm, convNum, maxChooseNum, outputNum, err := input()
	// if err != nil {
	// 	return
	// }
	// // algorithm, convNum, maxChooseNum, outputNum := 2, 3, 5, 2

	// if algorithm == 3 {
	// 	compare(convNum, maxChooseNum, outputNum)
	// } else {
	// 	p := process.NewProcessor(convNum, maxChooseNum, algorithm-1)
	// 	p.Prepare()
	// 	p.LoadData()
	// 	p.Match()
	// 	p.Output(process.OutputType(outputNum))
	// }

	// pause()

	eval.Evaluation(26, 8)
}

func compare(convNum, maxChooseNum, outputNum int) {
	ps := make([]process.MatchProcessor, 2)
	wg := sync.WaitGroup{}
	for i := 1; i >= 0; i-- {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			ps[i] = process.NewProcessor(convNum, maxChooseNum, i)
			ps[i].Prepare()
			ps[i].LoadData()
			ps[i].Match()
		}()
	}

	wg.Wait()

	fmt.Println("------ 生成对子后分配 ------")
	ps[0].Output(process.OutputType(outputNum))
	fmt.Println()

	fmt.Println("------ 逐轮生成对话 ------")
	ps[1].Output(process.OutputType(outputNum))
}

func input() (a, b, c, d int, err error) {
	// 读取第一个整数
	fmt.Println("请输入算法（输入数字后按下回车键）: ")
	fmt.Println("1. 生成对子后分配")
	fmt.Println("2. 逐轮生成对话")
	fmt.Println("3. 两种算法对比")
	if _, err = fmt.Scanln(&a); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()

	// 读取第一个整数
	fmt.Print("请输入对话轮数（输入数字后按下回车键）: ")
	if _, err = fmt.Scanln(&b); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()

	// 读取第二个整数
	fmt.Print("请输入最大可选择人数（输入数字后按下回车键）: ")
	if _, err = fmt.Scanln(&c); err != nil {
		fmt.Println("无效的输入:", err)
		return
	}
	fmt.Println()

	// 读取第三个整数
	fmt.Println("请选择数据输出模式（输入数字后按下回车键）")
	fmt.Println("1. 输出桌号")
	fmt.Println("2. 输出匹配类型")
	if _, err = fmt.Scanln(&d); err != nil {
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
