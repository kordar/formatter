package formatter

import (
	"fmt"
	"github.com/spf13/cast"
	"testing"
	"time"
)

type Demo struct {
	A time.Time `formatter:"time:2006-01-02"`
	B string    `json:"test" formatter:"roundstr:2"`
}

func TestAA(t *testing.T) {
	demos := []*Demo{
		{A: time.Now(), B: "123"},
		{A: time.Now(), B: "456"},
	}
	view := ToDataList(&demos, nil)
	fmt.Println(view)
}

func TestA22(t *testing.T) {
	result, err := cast.ToSliceE([]interface{}{1, 2, 3})
	fmt.Println(result, err) // 输出：[1 2 3]
}
