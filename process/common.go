package process

func isBetterConv(c1, c2 []Conversation) bool {
	t1, two1, rand1 := getTot(c1)
	t2, two2, rand2 := getTot(c2)
	if t1 == t2 {
		if two1 == two2 {
			return rand1 < rand2
		}
		return two1 > two2
	}
	return t1 > t2
}

func getTot(c []Conversation) (tot, twoTot, randTot int) {
	for _, conv := range c {
		tot += len(conv)
		for _, pair := range conv {
			if pair.MatchType == twoWayType {
				twoTot++
			} else if pair.MatchType == randType {
				randTot++
			}
		}
	}
	return
}
