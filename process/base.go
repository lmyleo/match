package process

import (
	"fmt"
	"match/entry"
	"match/util"
	"math"
	"sort"
	"strconv"

	"github.com/pkg/errors"
)

type BaseProcessor struct {
	ConvNum      int // 对话轮数
	MaxChooseNum int // 选择人数
	Entries      map[int64]*entry.Entry
	seq          []int64

	Convs     []Conversation // 多轮对话
	AllPairs  Conversation
	entryLeft map[int64]int
	entryStat map[int64]map[int]bool // 实体已加入的对话
	finish    bool
}

type Conversation []*Pair

type Pair struct {
	IDs       [2]int64
	MatchType string
}

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
	p.entryLeft = make(map[int64]int)
	p.entryStat = make(map[int64]map[int]bool)
	for i := 0; i < p.ConvNum; i++ {
		conv := Conversation{}
		p.Convs = append(p.Convs, conv)
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

	// 初始化状态
	for i := 0; i < p.ConvNum; i++ {
		for _, id := range p.seq {
			if i == 0 {
				p.entryLeft[id] = p.ConvNum
				p.entryStat[id] = make(map[int]bool)
			}
			p.entryStat[id][i] = false
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
		convWithType := make(map[string]Conversation)
		for _, pair := range conv {
			if _, ok := convWithType[pair.MatchType]; !ok {
				convWithType[pair.MatchType] = Conversation{}
			}
			convWithType[pair.MatchType] = append(convWithType[pair.MatchType], pair)
		}
		fmt.Printf("第 %d 轮对话\n", i+1)
		for _, matchType := range matchTypes {
			fmt.Printf("%s: ", matchType)
			for _, pair := range convWithType[matchType] {
				// fmt.Printf("%d桌%v  ", total+1, pair)
				fmt.Printf("[%d %d] ", pair.IDs[0], pair.IDs[1])
				total++
			}
			fmt.Println()

			switch matchType {
			case twoWayType:
				twoWayTot += len(convWithType[matchType])
			case oneWayType:
				oneWayTot += len(convWithType[matchType])
			case randType:
				randTot += len(convWithType[matchType])
			}
		}
		fmt.Printf("总对话数: %d\n", total)
		fmt.Println()
	}
	fmt.Printf("对话数：双向-%d；单向-%d；随机-%d\n", twoWayTot, oneWayTot, randTot)
	return nil
}

func (p *BaseProcessor) OutputData() error {
	for i, conv := range p.Convs {
		total := 0
		fmt.Printf("第 %d 轮对话\n", i+1)
		for _, pair := range conv {
			if total > 0 && total%5 == 0 {
				fmt.Println()
			}
			fmt.Printf("%d桌%v  ", total+1, pair.IDs)
			// fmt.Printf("%v  ", pair)
			total++
		}
		// fmt.Println()
	}
	fmt.Println()
	return nil
}

func (p *BaseProcessor) Match() {
	// 生成双向匹配
	p.genTwoWay()
	// 生成单向匹配
	p.genOneWay()
	// 生成随机匹配
	p.genRandom()
	// 将所有匹配安排到对话中
	p.distributeDfs(0)

	return
}

// genTwoWay 生成双向匹配
func (p *BaseProcessor) genTwoWay() {
	pairs := make([]*Pair, 0)
	set := make(map[string]struct{}, 0)
	for _, id1 := range p.seq {
		for _, id2 := range p.Entries[id1].GetChooseList() {
			if _, ok := set[fmt.Sprintf("%d_%d", id2, id1)]; !ok && p.Entries[id2].Find(id1) {
				set[fmt.Sprintf("%d_%d", id1, id2)] = struct{}{}
				p.entryLeft[id1]--
				p.entryLeft[id2]--
				pairs = append(pairs, &Pair{
					IDs:       [2]int64{id1, id2},
					MatchType: twoWayType,
				})
			}
		}
	}
	pairs = p.filterPair(pairs)
	for _, pair := range pairs {
		p.matchOnePair(pair)
	}
	fmt.Printf("\n双向总数：%d\n", len(pairs))
}

// genOneWay 生成单向匹配
func (p *BaseProcessor) genOneWay() {
	var pairs []*Pair
	// 生成匹配
	set := make(map[string]struct{}, 0)
	for _, id1 := range p.seq {
		for _, id2 := range p.Entries[id1].GetChooseList() {
			if _, ok := set[fmt.Sprintf("%d_%d", id2, id1)]; !ok && !p.Entries[id2].Find(id1) {
				set[fmt.Sprintf("%d_%d", id1, id2)] = struct{}{}
				p.entryLeft[id1]--
				p.entryLeft[id2]--
				pairs = append(pairs, &Pair{
					IDs:       [2]int64{id1, id2},
					MatchType: oneWayType,
				})
			}
		}
	}
	pairs = p.filterPair(pairs)
	for _, pair := range pairs {
		p.matchOnePair(pair)
	}
	fmt.Printf("\n单向总数：%d\n", len(pairs))
}

// genOneWay 生成单向匹配
func (p *BaseProcessor) genRandom() {
	var pairs []*Pair
	ids := make([]int64, 0)
	for id, left := range p.entryLeft {
		if left > 0 {
			ids = append(ids, id)
		}
	}

	// 生成匹配
	set := make(map[string]struct{}, 0)
	for _, id1 := range ids {
		for _, id2 := range ids {
			if _, ok := set[fmt.Sprintf("%d_%d", id2, id1)]; !ok && id1 != id2 && !p.isMatched(id1, id2) {
				set[fmt.Sprintf("%d_%d", id1, id2)] = struct{}{}
				p.entryLeft[id1]--
				p.entryLeft[id2]--
				pairs = append(pairs, &Pair{
					IDs:       [2]int64{id1, id2},
					MatchType: randType,
				})
			}
		}
	}
	pairs = p.filterPair(pairs)
	for _, pair := range pairs {
		p.matchOnePair(pair)
	}
	fmt.Printf("\n随机总数：%d\n", len(pairs))
}

func (p *BaseProcessor) distributeDfs(num int) {
	if num == len(p.AllPairs) {
		p.finish = true
		return
	}
	// fmt.Printf("num: %d\n", num)
	id1, id2 := p.AllPairs[num].IDs[0], p.AllPairs[num].IDs[1]
	for pos := 0; pos < p.ConvNum; pos++ {
		if p.entryStat[id1][pos] || p.entryStat[id2][pos] || len(p.Convs[pos]) >= len(p.seq)/2 {
			continue
		}
		// 将 id1, id2 安排到此轮对话
		p.Convs[pos] = append(p.Convs[pos], p.AllPairs[num])
		p.entryStat[id1][pos] = true
		p.entryStat[id2][pos] = true
		p.distributeDfs(num + 1)
		if p.finish {
			return
		}
		p.Convs[pos] = p.Convs[pos][:len(p.Convs[pos])-1]
		p.entryStat[id1][pos] = false
		p.entryStat[id2][pos] = false
	}
	return
}

// filterPair 过滤多余匹配（涉及状态：p.EntryLeft）
func (p *BaseProcessor) filterPair(pairs []*Pair) (newPairs []*Pair) {
	mp := make(map[int64][]int64)

	for _, p := range pairs {
		id1, id2 := p.IDs[0], p.IDs[1]
		if _, ok := mp[id1]; !ok {
			mp[id1] = make([]int64, 0)
		}
		if _, ok := mp[id2]; !ok {
			mp[id2] = make([]int64, 0)
		}
		mp[id1] = append(mp[id1], id2)
		mp[id2] = append(mp[id2], id1)
	}

	for id, matches := range mp {
		if p.entryLeft[id] >= 0 {
			continue
		}
		sort.Slice(matches, func(i, j int) bool {
			id1, id2 := matches[i], matches[j]
			return p.entryLeft[id1] < p.entryLeft[id2]
		})
		// 过滤匹配
		filterCnt := -p.entryLeft[id]
		for i := 0; i < filterCnt && i < len(mp[id]); i++ {
			filterID := mp[id][i]
			p.entryLeft[filterID]++
			mp[filterID] = util.DeleteInt64(mp[filterID], id)
		}
		mp[id] = mp[id][filterCnt:]
		p.entryLeft[id] = 0
	}

	// 构造新 paires
	idSet := make(map[int64]map[int64]struct{})
	for id, matches := range mp {
		if _, ok := idSet[id]; !ok {
			idSet[id] = make(map[int64]struct{})
		}
		for _, id2 := range matches {
			idSet[id][id2] = struct{}{}
		}
	}
	for _, p := range pairs {
		id1, id2 := p.IDs[0], p.IDs[1]
		if _, ok := idSet[id1][id2]; ok {
			newPairs = append(newPairs, p)
		}
	}

	return newPairs
}

// matchOnePair 生成一对
func (p *BaseProcessor) matchOnePair(pair *Pair) bool {
	id1, id2 := pair.IDs[0], pair.IDs[1]
	if p.isMatched(id1, id2) {
		return false
	}
	p.Entries[id1].Match(id2)
	p.Entries[id2].Match(id1)

	fmt.Printf("[%d %d] ", id1, id2)
	p.AllPairs = append(p.AllPairs, pair)

	return true
}

func (p *BaseProcessor) addEntry(id int64, chooseList []int64) {
	p.Entries[id] = entry.NewEntry(id, chooseList)
	p.seq = append(p.seq, id)
}

func (p *BaseProcessor) isMatched(id1, id2 int64) bool {
	return p.Entries[id1].IsMatched(id2)
}

// func (p *BaseProcessor) reduceTwoWayMatch() {
// 	for e1 := range p.twoWayMatch {
// 		for len(p.twoWayMatch[e1]) > p.ConvNum {
// 			// 删除匹配数最多的人
// 			sort.Slice(p.twoWayMatch[e1], func(i, j int) bool {
// 				e2, e3 := p.twoWayMatch[e1][i], p.twoWayMatch[e1][j]
// 				return len(p.twoWayMatch[e2]) > len(p.twoWayMatch[e3]) // 降序排序
// 			})
// 			p.deleteTwoWayMatch(e1, p.twoWayMatch[e1][0])
// 		}
// 	}
// }

// func (p *BaseProcessor) addTwoWayConv() {
// 	// 按照双向匹配数排序，优先安排双向匹配多的人
// 	seq := make([]int64, 0)
// 	for e := range p.twoWayMatch {
// 		seq = append(seq, e)
// 	}
// 	sort.Slice(seq, func(i, j int) bool {
// 		e1, e2 := seq[i], seq[j]
// 		return len(p.twoWayMatch[e1]) > len(p.twoWayMatch[e2]) // 降序排序
// 	})
// 	// 将匹配加入对话
// 	for _, e1 := range seq {
// 		chooseList := p.twoWayMatch[e1]
// 		for _, e2 := range chooseList {
// 			if p.isMatched(e1, e2) {
// 				continue
// 			}
// 			for j := 0; j < p.ConvNum; j++ {
// 				if p.matchOnePair(twoWayType, j, e1, e2) {
// 					break
// 				}
// 			}
// 		}
// 	}
// }

// func (p *BaseProcessor) addTwoWayMatch(e1, e2 int64) {
// 	_, ok := p.twoWayMatch[e1]
// 	if !ok {
// 		p.twoWayMatch[e1] = make([]int64, 0)
// 	}
// 	_, ok = p.twoWayMatch[e2]
// 	if !ok {
// 		p.twoWayMatch[e2] = make([]int64, 0)
// 	}
// 	p.twoWayMatch[e1] = append(p.twoWayMatch[e1], e2)
// 	p.twoWayMatch[e2] = append(p.twoWayMatch[e2], e1)
// }

// func (p *BaseProcessor) deleteTwoWayMatch(e1, e2 int64) {
// 	p.twoWayMatch[e1] = util.DeleteInt64(p.twoWayMatch[e1], e2)
// 	p.twoWayMatch[e2] = util.DeleteInt64(p.twoWayMatch[e2], e1)
// }

// func (p *BaseProcessor) MatchMultiTimes(times int) {
// 	for i := 0; i < times; i++ {

// 	}
// }

// MatchOneConv 匹配一轮对话
// func (p *BaseProcessor) MatchOneConv() error {
// 	// 为保证公平，随机打乱顺序
// 	util.Shuffle(p.seq)

// 	p.TwoWayMatch() // 双向匹配

// 	p.OneWayMatch() // 单向匹配

// 	p.RandMatch() // 随机匹配

// 	return nil
// }

// TwoWayMatch 双向匹配
// func (p *BaseProcessor) TwoWayMatch() {
// 	for _, id := range p.seq {
// 		targets := p.getMatchableIDs(id) // 获取可匹配对象

//			for _, targetID := range targets {
//				target := p.Entries[targetID]
//				// 是否满足双向匹配
//				if target.Find(id) {
//					p.MatchOnePair(twoWayType, id, targetID)
//					break
//				}
//			}
//		}
//		return
//	}
//
// OneWayMatch 单向匹配
// func (p *BaseProcessor) OneWayMatch() {
// 	for _, id := range p.seq {
// 		targets := p.getMatchableIDs(id) // 获取可匹配对象

// 		for _, targetID := range targets {
// 			p.MatchOnePair(oneWayType, id, targetID)
// 			break
// 		}
// 	}
// 	return
// }

// RandMatch 随机匹配
// func (p *BaseProcessor) RandMatch() {
// 	unMatch := make([]int64, 0)
// 	// 找出所有未匹配的实体
// 	for _, id := range p.seq {
// 		if _, ok := p.curEntryStat[id]; !ok {
// 			unMatch = append(unMatch, id)
// 		}
// 	}
// 	p.randMatch[p.curConv] = unMatch
// 	// 随机匹配
// 	for i := 0; i < 100; i++ {
// 		pairNum := 0
// 		for i := 0; i+1 < len(unMatch); i += 2 {
// 			e1, e2 := unMatch[i], unMatch[i+1]
// 			if !p.isMatched(e1, e2) {
// 				pairNum++
// 				continue // e1 与 e2 匹配，查看下两个实体
// 			}
// 			// e1 与 e2 无法匹配，需要 e2 与 e3 调换位置
// 			swap := false
// 			for j := len(unMatch) - 1; j >= 0; j-- {
// 				e3 := unMatch[j]
// 				if j == i || j == i+1 || p.isMatched(e1, e3) {
// 					continue
// 				}
// 				// 查看换位后是否会破坏已有匹配
// 				if j < i {
// 					var e4 int64 // e3 之前的匹配者
// 					if j%2 == 0 {
// 						e4 = unMatch[j+1]
// 					} else {
// 						e4 = unMatch[j-1]
// 					}
// 					if p.isMatched(e2, e4) {
// 						continue // e2 与 e4 无法匹配，不能换位
// 					}
// 				}
// 				// 换位
// 				swap = true
// 				unMatch[i+1], unMatch[j] = unMatch[j], unMatch[i+1]
// 				pairNum++
// 				break
// 			}
// 			// 无法匹配对话
// 			if !swap {
// 				break
// 			}
// 		}
// 		if pairNum == len(unMatch)/2 {
// 			break
// 		}
// 		util.Shuffle(unMatch)
// 	}

// 	// 生成匹配
// 	for i := 0; i+1 < len(unMatch); i += 2 {
// 		p.MatchOnePair(randType, unMatch[i], unMatch[i+1])
// 	}
// 	return
// }

// func (p *BaseProcessor) getMatchableIDs(id int64) (matchableIDs []int64) {
// 	if p.curEntryStat[id] {
// 		return
// 	}
// 	for target := range p.Entries[id].GetChooseSet() {
// 		// target 在此轮未匹配 && target 未曾与 id 匹配
// 		if !p.curEntryStat[target] && !p.Entries[target].IsMatched(id) {
// 			matchableIDs = append(matchableIDs, target)
// 		}
// 	}
// 	return
// }

// func (p *BaseProcessor) addConv(conv int, pair *Pair) {
// 	id1, id2 := pair.IDs[0], pair.IDs[1]
// 	p.entryStat[id1][conv] = true
// 	p.entryStat[id2][conv] = true
// }

// func (p *BaseProcessor) canAdd(conv int, pair *Pair) bool {
// 	id1, id2 := pair.IDs[0], pair.IDs[1]
// 	return !p.entryStat[id1][conv] && !p.entryStat[id2][conv]
// }
