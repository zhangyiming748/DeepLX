package translate

import (
	"log"
	"testing"
)

func TestTranslateByDeepLX(t *testing.T) {
	source := "auto"
	target := "zh"
	text := "hello"
	lx, err := TranslateByDeepLX(source, target, text, "html", "", "")
	if err != nil {
		log.Fatalln(err)
	}
	t.Log(lx)
}
