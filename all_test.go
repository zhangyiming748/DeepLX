package DeepLx

import "testing"

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
