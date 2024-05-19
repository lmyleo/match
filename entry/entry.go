package entry

type Entry struct {
	id         int64
	property   string
	chooseList []int64
	chooseSet  map[int64]struct{} // 1: 未匹配; 2: 已匹配;
	matchSet   map[int64]struct{}
}

type Status int

const (
	UnMatch Status = 1
	Match   Status = 2
)

func NewEntry(id int64, chooseList []int64) *Entry {
	e := &Entry{chooseList: chooseList, chooseSet: make(map[int64]struct{}), matchSet: make(map[int64]struct{})}
	e.id = id
	for _, id := range e.chooseList {
		if _, ok := e.chooseSet[id]; ok || id == e.id {
			continue
		}
		e.chooseSet[id] = struct{}{}
	}
	return e
}
func (e *Entry) GetChooseList() []int64 {
	return e.chooseList
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
