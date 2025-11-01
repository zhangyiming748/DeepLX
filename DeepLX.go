package DeepLX

import (
	"log"
	"time"

	"github.com/zhangyiming748/DeepLX/translate"
)

func DeepLX(src string) (dst string) {
	result, err := translate.TranslateByDeepLX("auto", "zh", src, "", "", "")
	// 无论如何都打印两个变量=
	if err != nil {
		log.Fatalf("%v", err)
	}
	if result.Data == "" {
		time.Sleep(time.Second)
		return DeepLX(src)
	}

	return result.Data
}
