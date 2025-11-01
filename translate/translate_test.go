package translate

import (
	"log"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	src := "エロ"
	dst := Trans(src)
	log.Println(dst)
}

func Trans(src string) (dst string) {
	result, err := TranslateByDeepLX("jp", "zh", src, "", "", "")
	// 无论如何都打印两个变量=
	if err != nil {
		log.Fatalf("%v", err)
	}
	if result.Data == "" {
		time.Sleep(time.Second)
		return Trans(src)
	} else {
		return result.Data
	}
}
