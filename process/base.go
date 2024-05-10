package process

import (
	"match/entry"
	"match/util"
)

type BaseProcessor struct {
	ConvNum   int
	ChooseNum int
	Entries   map[int64]*entry.Entry
	seq       []int64
	Convs     []Conversation
}

type Conversation []Pair

type Pair [2]int

func (p *BaseProcessor) LoadData() error {
	return nil
}

func (p *BaseProcessor) OutputData() error {
	return nil
}

func (p *BaseProcessor) Match() error {
	for i := 0; i < p.ConvNum; i++ {
		p.MatchOneConv()
	}
	return nil
}

// MatchOneConv 匹配一轮对话
func (p *BaseProcessor) MatchOneConv() error {
	// 为保证公平，随机打乱顺序
	util.Shuffle(p.seq)

	// 双向匹配

	// 单向匹配

	return nil
}

// TwoWayMatch 双向匹配
func (p *BaseProcessor) TwoWayMatch() error {
	for _, id := range p.seq {
		e := p.Entries[id]
		chooseSet := e.GetChooseSet()

		for targetID, stat := range chooseSet {
			// 查找 e 是否在其未匹配目标的选择中
			if stat == entry.UnMatch {
				target := p.Entries[targetID]
				if target.Find(e.GetId()) {
					p.MatchOnePair(id, targetID)
				}
			}
		}
	}
	return nil
}

// MatchOnePair 匹配一对对话
func (p *BaseProcessor) MatchOnePair(id1, id2 int64) {
	// 改变实体状态
	e1, e2 := p.Entries[id1], p.Entries[id2]
	e1.Match(id2)
	e2.Match(id1)
	// TODO 将这一对加入此轮对话

}
