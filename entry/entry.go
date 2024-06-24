package entry

type Entry struct {
	id         int64
	property   string
	chooseList []int64
	chooseSeq  map[int64]int64    // 喜好顺序
	matchSet   map[int64]struct{} // 匹配的人的分数
}

type Status int

const (
	UnMatch Status = 1
	Match   Status = 2
)

func NewEntry(id int64, chooseList []int64) *Entry {
	e := &Entry{chooseList: make([]int64, 0), chooseSeq: make(map[int64]int64), matchSet: make(map[int64]struct{})}
	e.id = id
	seq := int64(0)
	for _, id := range chooseList {
		if _, ok := e.chooseSeq[id]; ok || id == e.id {
			continue
		}
		e.chooseList = append(e.chooseList, id)
		e.chooseSeq[id] = seq
		seq++
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
