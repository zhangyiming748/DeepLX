package DeepLx

import (
	"fmt"
	"testing"
)

func TestUsage(t *testing.T) {
	source := "EN"
	target := "ZH"
	text := "hello"
	lx, err := TranslateByDeepLX(source, target, text, "")
	if err != nil {
		return
	} else {
		t.Log(lx)
	}
}
func TestString(t *testing.T) {
	fmt.Println("您好")
	fmt.Println("喂")
	fmt.Println("您好！请问您是哪位？")
	fmt.Println("你好")
	fmt.Println("")
}
