package process

type MatchProcessor interface {
	LoadData() error
	Match() error
	OutputData() error
}
