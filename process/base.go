package process

import (
	"fmt"
	"match/entry"
	"match/util"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type BaseProcessor struct {
	ConvNum      int // 对话轮数
	MaxChooseNum int // 选择人数
	Entries      map[int64]*entry.Entry
	seq          []int64

	Convs        []Conversation // 多轮对话
	curConv      int            //  当前对话轮数
	curEntryStat map[int64]bool // 当前对话轮数中，实体是否已经匹配
	randMatch    [][]int64      // 每一轮中随机匹配的实体
}

type Conversation map[string][]Pair

type Pair [2]int64

const (
	twoWayType = "双向匹配"
	oneWayType = "单向匹配"
	randType   = "随机匹配"
)

var matchTypes = []string{twoWayType, oneWayType, randType}

func NewBaseProcessor(convNum, maxChooseNum int) *BaseProcessor {
	p := &BaseProcessor{ConvNum: convNum, MaxChooseNum: maxChooseNum}
	return p
}

func (p *BaseProcessor) Prepare() {
	p.Entries = make(map[int64]*entry.Entry)
	for i := 0; i < p.ConvNum; i++ {
		conv := Conversation{}
		p.Convs = append(p.Convs, conv)
		ids := make([]int64, 0)
		p.randMatch = append(p.randMatch, ids)
	}
}

func (p *BaseProcessor) LoadData() error {
	inputRows, err := util.ReadExcel()
	if err != nil {
		return errors.Wrap(err, "read input")
	}
	for row, content := range inputRows {
		if row == 0 {
			continue
		}
		// 读取每个人的 id & 其选择的多个 id
		var id int64
		choose := make([]int64, 0)
		for col, value := range content {
			if col == 1 {
				id, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					return errors.Wrapf(err, "conv id to int faild value[%s]", value)
				}
			}
			if col == 2 {
				choose = util.GetNumbers(value)
				p.addEntry(id, choose[0:int(math.Min(float64(p.MaxChooseNum), float64(len(choose))))])
			}
		}
	}
	return nil
}

func (p *BaseProcessor) Output(outputType OutputType) {
	if outputType == ForUser {
		p.OutputData()
	} else if outputType == ForDev {
		p.OutputDebug()
	}
}

func (p *BaseProcessor) OutputDebug() error {
	twoWayTot, oneWayTot, randTot := 0, 0, 0
	for i, conv := range p.Convs {
		total := 0
		fmt.Printf("第 %d 轮对话\n", i+1)
		for _, matchType := range matchTypes {
			fmt.Printf("%s: ", matchType)
			for _, pair := range conv[matchType] {
				// fmt.Printf("%d桌%v  ", total+1, pair)
				fmt.Printf("%v  ", pair)
				total++
			}
			fmt.Println()

			switch matchType {
			case twoWayType:
				twoWayTot += len(conv[matchType])
			case oneWayType:
				oneWayTot += len(conv[matchType])
			case randType:
				randTot += len(conv[matchType])
			}
		}
		fmt.Printf("总对话数: %d\n", total)
		if total < len(p.seq)/2 {
			fmt.Printf("随机匹配实体：%v\n", p.randMatch[i])
		}
		fmt.Println()
	}
	fmt.Printf("对话数：双向-%d；单向-%d；随机-%d\n", twoWayTot, oneWayTot, randTot)
	return nil
}

func (p *BaseProcessor) OutputData() error {
	for i, conv := range p.Convs {
		total := 0
		fmt.Printf("第 %d 轮对话\n", i+1)
		for _, matchType := range matchTypes {
			for _, pair := range conv[matchType] {
				if total > 0 && total%5 == 0 {
					fmt.Println()
				}
				fmt.Printf("%d桌%v  ", total+1, pair)
				// fmt.Printf("%v  ", pair)
				total++
			}
			// fmt.Println()
		}
		fmt.Printf("\n总对话数: %d\n", total)
		if total < len(p.seq)/2 {
			fmt.Printf("随机匹配实体：%v\n", p.randMatch[i])
		}
		fmt.Println()
	}
	return nil
}

func (p *BaseProcessor) Match() {
	// 生成双向匹配

	// 按照不重复原则分配到每轮对话中

	// 生成单向匹配

	return
}

func (p *BaseProcessor) MatchMultiTimes(times int) {
	for i := 0; i < times; i++ {

	}
}

// MatchOneConv 匹配一轮对话
func (p *BaseProcessor) MatchOneConv() error {
	// 为保证公平，随机打乱顺序
	util.Shuffle(p.seq)

	p.TwoWayMatch() // 双向匹配

	p.OneWayMatch() // 单向匹配

	p.RandMatch() // 随机匹配

	return nil
}

// TwoWayMatch 双向匹配
func (p *BaseProcessor) TwoWayMatch() {
	for _, id := range p.seq {
		targets := p.getMatchableIDs(id) // 获取可匹配对象

		for _, targetID := range targets {
			target := p.Entries[targetID]
			// 是否满足双向匹配
			if target.Find(id) {
				p.MatchOnePair(twoWayType, id, targetID)
				break
			}
		}
	}
	return
}

// OneWayMatch 单向匹配
func (p *BaseProcessor) OneWayMatch() {
	for _, id := range p.seq {
		targets := p.getMatchableIDs(id) // 获取可匹配对象

		for _, targetID := range targets {
			p.MatchOnePair(oneWayType, id, targetID)
			break
		}
	}
	return
}

// RandMatch 随机匹配
func (p *BaseProcessor) RandMatch() {
	unMatch := make([]int64, 0)
	// 找出所有未匹配的实体
	for _, id := range p.seq {
		if _, ok := p.curEntryStat[id]; !ok {
			unMatch = append(unMatch, id)
		}
	}
	p.randMatch[p.curConv] = unMatch
	// 随机匹配
	for i := 0; i < 100; i++ {
		pairNum := 0
		for i := 0; i+1 < len(unMatch); i += 2 {
			e1, e2 := unMatch[i], unMatch[i+1]
			if !p.isMatched(e1, e2) {
				pairNum++
				continue // e1 与 e2 匹配，查看下两个实体
			}
			// e1 与 e2 无法匹配，需要 e2 与 e3 调换位置
			swap := false
			for j := len(unMatch) - 1; j >= 0; j-- {
				e3 := unMatch[j]
				if j == i || j == i+1 || p.isMatched(e1, e3) {
					continue
				}
				// 查看换位后是否会破坏已有匹配
				if j < i {
					var e4 int64 // e3 之前的匹配者
					if j%2 == 0 {
						e4 = unMatch[j+1]
					} else {
						e4 = unMatch[j-1]
					}
					if p.isMatched(e2, e4) {
						continue // e2 与 e4 无法匹配，不能换位
					}
				}
				// 换位
				swap = true
				unMatch[i+1], unMatch[j] = unMatch[j], unMatch[i+1]
				pairNum++
				break
			}
			// 无法匹配对话
			if !swap {
				break
			}
		}
		if pairNum == len(unMatch)/2 {
			break
		}
		util.Shuffle(unMatch)
	}

	// 生成匹配
	for i := 0; i+1 < len(unMatch); i += 2 {
		p.MatchOnePair(randType, unMatch[i], unMatch[i+1])
	}
	return
}

// MatchOnePair 匹配一对对话
func (p *BaseProcessor) MatchOnePair(matchType string, id1, id2 int64) {
	// 判断是否能匹配
	if p.curEntryStat[id1] || p.curEntryStat[id2] || p.Entries[id1].IsMatched(id2) {
		return
	}

	// 改变实体状态
	e1, e2 := p.Entries[id1], p.Entries[id2]
	e1.Match(id2)
	e2.Match(id1)

	// 将这一对加入此轮对话
	p.Convs[p.curConv][matchType] = append(p.Convs[p.curConv][matchType], Pair{id1, id2})
	p.curEntryStat[id1] = true
	p.curEntryStat[id2] = true
}

func (p *BaseProcessor) getMatchableIDs(id int64) (matchableIDs []int64) {
	if p.curEntryStat[id] {
		return
	}
	for target := range p.Entries[id].GetChooseSet() {
		// target 在此轮未匹配 && target 未曾与 id 匹配
		if !p.curEntryStat[target] && !p.Entries[target].IsMatched(id) {
			matchableIDs = append(matchableIDs, target)
		}
	}
	return
}

func (p *BaseProcessor) addEntry(id int64, chooseList []int64) {
	p.Entries[id] = entry.NewEntry(id, chooseList)
	p.seq = append(p.seq, id)
}

func (p *BaseProcessor) isMatched(e1, e2 int64) bool {
	return p.Entries[e1].IsMatched(e2)
}
