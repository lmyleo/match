package entry

func (e *Entry) Find(id int64) bool {
	stat, ok := e.chooseSet[id]
	if ok && stat == UnMatch {
		return true
	}
	return false
}

func (e *Entry) Match(id int64) {
	e.chooseSet[id] = Match
}

func (e *Entry) IsMatched(id int64) bool {
	stat, ok := e.chooseSet[id]
	if ok && stat == Match {
		return true
	}
	return false
}
