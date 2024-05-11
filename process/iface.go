package process

type OutputType int

const (
	ForUser OutputType = 1
	ForDev  OutputType = 2
)

type MatchProcessor interface {
	Prepare()
	LoadData() error
	Match()
	Output(outputType OutputType)
}

func NewProcessor(convNum, maxChooseNum int) MatchProcessor {
	return NewBaseProcessor(convNum, maxChooseNum)
}
