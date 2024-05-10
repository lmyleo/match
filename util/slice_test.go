package util_test

import (
	"fmt"
	"match/util"
	"testing"

	. "github.com/bytedance/mockey"
)

func Test_ShuffleEntry(t *testing.T) {
	PatchConvey("Test_ShuffleEntry", t, func() {
		PatchConvey("case 1", func() {
			entries := make([]int64, 0)
			for i := 0; i < 100; i++ {
				entries = append(entries, int64(i))
			}
			for i := 0; i < 100; i++ {
				// So(string(entries[i].GetId()), ShouldEqual, strconv.Itoa(i))
				fmt.Println(entries[i])
			}
			fmt.Println("-------------- ShuffleEntry -------------")
			util.Shuffle(entries)
			for i := 0; i < 100; i++ {
				// So(string(entries[i].GetId()), ShouldEqual, strconv.Itoa(i))
				fmt.Println(entries[i])
			}
		})
	})
}
