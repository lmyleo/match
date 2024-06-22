package process

type OutputType int

const (
	ForUser OutputType = 1
	ForDev  OutputType = 2

	twoWayType = "双向匹配"
	oneWayType = "单向匹配"
	randType   = "随机匹配"

	Base = 0
	Seq  = 1
)

type Conversation []*Pair

type Pair struct {
	IDs       [2]int64
	MatchType string
	Score     int64
}

type MatchProcessor interface {
	Prepare()
	LoadData() error
	Match()
	Output(outputType OutputType)
}

func NewProcessor(convNum, maxChooseNum, processType int) MatchProcessor {
	if processType == Base {
		return NewBaseProcessor(convNum, maxChooseNum)
	} else if processType == Seq {
		return NewSeqProcessor(convNum, maxChooseNum)
	}
	return nil
}
