package entry

type Entry struct {
	id         int64
	property   string
	chooseList []int64
	chooseSet  map[int64]Status // 1: 未匹配; 2: 已匹配;
}

type Status int

const (
	UnMatch Status = 1
	Match   Status = 2
)

func (e *Entry) GetChooseList() []int64 {
	return e.chooseList
}

func (e *Entry) GetChooseSet() map[int64]Status {
	return e.chooseSet
}

func (e *Entry) GetId() int64 {
	return e.id
}

func (e *Entry) SetId(id int64) {
	e.id = id
	return
}

func (e *Entry) GetProperty() string {
	return e.property
}

func (e *Entry) SetChooseSet(set map[int64]Status) {
	e.chooseSet = set
}

func (e *Entry) SetChooseList(list []int64) {
	e.chooseList = list
}
