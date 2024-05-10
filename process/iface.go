package process

type MatchProcessor interface {
	Prepare()
	LoadData() error
	Match() error
	OutputData() error
}
