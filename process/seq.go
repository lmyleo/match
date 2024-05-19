package process

import (
	"match/util"
)

type SeqProcessor struct {
	BaseProcessor
	curConv   int
	randMatch [][]int64
}

func NewSeqProcessor(convNum, maxChooseNum int) *SeqProcessor {
	p := &SeqProcessor{BaseProcessor: BaseProcessor{ConvNum: convNum, MaxChooseNum: maxChooseNum}}
	return p
}

func (p *SeqProcessor) Prepare() {
	p.BaseProcessor.Prepare()
	p.randMatch = make([][]int64, p.ConvNum)
	for i := 0; i < p.ConvNum; i++ {
		p.randMatch[i] = make([]int64, 0)
	}
}

func (p *SeqProcessor) Match() {
	p.MatchMultiTimes(10)
}

// MatchConvs 匹配对话
func (p *SeqProcessor) MatchConvs() error {
	// 为保证公平，随机打乱顺序
	util.Shuffle(p.seq)

	for i := 0; i < p.ConvNum; i++ {
		p.curConv = i

		p.TwoWayMatch() // 双向匹配

		p.OneWayMatch() // 单向匹配

		p.RandMatch() // 随机匹配
	}

	return nil
}

// TwoWayMatch 双向匹配
func (p *SeqProcessor) TwoWayMatch() {
	for _, id := range p.seq {
		targets := p.getMatchableIDs(id) // 获取可匹配对象

		for _, targetID := range targets {
			target := p.Entries[targetID]
			// 是否满足双向匹配
			if target.Find(id) {
				p.matchOnePair(twoWayType, id, targetID)
				break
			}
		}
	}
	return
}

// OneWayMatch 单向匹配
func (p *SeqProcessor) OneWayMatch() {
	for _, id := range p.seq {
		targets := p.getMatchableIDs(id) // 获取可匹配对象

		for _, targetID := range targets {
			p.matchOnePair(oneWayType, id, targetID)
			break
		}
	}
	return
}

// RandMatch 随机匹配
func (p *SeqProcessor) RandMatch() {
	unMatch := make([]int64, 0)
	// 找出所有未匹配的实体
	for _, id := range p.seq {
		if !p.entryStat[id][p.curConv] {
			unMatch = append(unMatch, id)
		}
	}
	p.randMatch[p.curConv] = unMatch
	// 随机匹配
	for ii := 0; ii < 100; ii++ {
		pairNum := 0
		for i := 0; i+1 < len(unMatch); i += 2 {
			id1, id2 := unMatch[i], unMatch[i+1]
			if !p.isMatched(id1, id2) {
				pairNum++
				continue // id1 与 id2 匹配，查看下两个实体
			}
			// id1 与 id2 无法匹配，需要 id2 与 id3 调换位置
			swap := false
			for j := len(unMatch) - 1; j >= 0; j-- {
				id3 := unMatch[j]
				if j == i || j == i+1 || p.isMatched(id1, id3) {
					continue
				}
				// 查看换位后是否会破坏已有匹配
				if j < i {
					var id4 int64 // id3 之前的匹配者
					if j%2 == 0 {
						id4 = unMatch[j+1]
					} else {
						id4 = unMatch[j-1]
					}
					if p.isMatched(id2, id4) {
						continue // id2 与 id4 无法匹配，不能换位
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
		p.matchOnePair(randType, unMatch[i], unMatch[i+1])
	}
	return
}

func (p *SeqProcessor) matchOnePair(matchType string, id1, id2 int64) bool {
	if p.isMatched(id1, id2) {
		return false
	}
	p.Entries[id1].Match(id2)
	p.Entries[id2].Match(id1)
	p.entryStat[id1][p.curConv] = true
	p.entryStat[id2][p.curConv] = true
	p.entryLeft[id1]--
	p.entryLeft[id2]--

	pair := &Pair{IDs: [2]int64{id1, id2}, MatchType: matchType}
	p.Convs[p.curConv] = append(p.Convs[p.curConv], pair)

	return true
}

func (p *SeqProcessor) getMatchableIDs(id int64) (matchableIDs []int64) {
	if p.entryStat[id][p.curConv] {
		return
	}
	for _, target := range p.Entries[id].GetChooseList() {
		// target 在此轮未匹配 && target 未曾与 id 匹配
		if !p.entryStat[target][p.curConv] && !p.Entries[target].IsMatched(id) {
			matchableIDs = append(matchableIDs, target)
		}
	}
	return
}
