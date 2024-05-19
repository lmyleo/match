package entry

func (e *Entry) Find(id int64) bool {
	_, ok := e.chooseSet[id]
	return ok
}

func (e *Entry) Match(id int64) {
	e.matchSet[id] = struct{}{}
}

func (e *Entry) IsMatched(id int64) bool {
	_, ok := e.matchSet[id]
	return ok
}

func (e *Entry) Clear() {
	e.matchSet = make(map[int64]struct{})
}
