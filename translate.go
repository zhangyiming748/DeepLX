package DeepLx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/andybalholm/brotli"
	"github.com/tidwall/gjson"
)

type req struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Result  struct {
		Texts []struct {
			Alternatives []struct {
				Text string `json:"text"`
			} `json:"alternatives"`
			Text string `json:"text"`
		} `json:"texts"`
		Lang              string `json:"lang"`
		LangIsConfident   bool   `json:"lang_is_confident"`
		DetectedLanguages struct {
			EN          float64 `json:"EN"`
			DE          float64 `json:"DE"`
			FR          float64 `json:"FR"`
			ES          float64 `json:"ES"`
			PT          float64 `json:"PT"`
			IT          float64 `json:"IT"`
			NL          float64 `json:"NL"`
			PL          float64 `json:"PL"`
			RU          float64 `json:"RU"`
			ZH          float64 `json:"ZH"`
			JA          float64 `json:"JA"`
			BG          float64 `json:"BG"`
			CS          float64 `json:"CS"`
			DA          float64 `json:"DA"`
			EL          float64 `json:"EL"`
			ET          float64 `json:"ET"`
			FI          float64 `json:"FI"`
			HU          float64 `json:"HU"`
			LT          float64 `json:"LT"`
			LV          float64 `json:"LV"`
			RO          float64 `json:"RO"`
			SK          float64 `json:"SK"`
			SL          float64 `json:"SL"`
			SV          float64 `json:"SV"`
			TR          float64 `json:"TR"`
			ID          float64 `json:"ID"`
			UK          float64 `json:"UK"`
			KO          float64 `json:"KO"`
			NB          float64 `json:"NB"`
			AR          float64 `json:"AR"`
			Unsupported float64 `json:"unsupported"`
		} `json:"detectedLanguages"`
	} `json:"result"`
}

func initDeepLXData(sourceLang string, targetLang string) *PostData {
	hasRegionalVariant := false
	targetLangParts := strings.Split(targetLang, "-")

	// targetLang can be "en", "pt", "pt-PT", "pt-BR"
	// targetLangCode is the first part of the targetLang, e.g. "pt" in "pt-PT"
	targetLangCode := targetLangParts[0]
	if len(targetLangParts) > 1 {
		hasRegionalVariant = true
	}

	commonJobParams := CommonJobParams{
		WasSpoken:    false,
		TranscribeAS: "",
	}
	if hasRegionalVariant {
		commonJobParams.RegionalVariant = targetLang
	}

	return &PostData{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_texts",
		Params: Params{
			Splitting: "newlines",
			Lang: Lang{
				SourceLangUserSelected: sourceLang,
				TargetLang:             targetLangCode,
			},
			CommonJobParams: commonJobParams,
		},
	}
}

func TranslateByDeepLX(sourceLang string, targetLang string, translateText string, proxyURL string) (string, error) {
	id := getRandomNumber()
	if sourceLang == "" {
		lang := whatlanggo.DetectLang(translateText)
		deepLLang := strings.ToUpper(lang.Iso6391())
		sourceLang = deepLLang
	}
	// If target language is not specified, set it to English
	if targetLang == "" {
		targetLang = "EN"
	}
	// Handling empty translation text
	if translateText == "" {
		return "", errors.New("No text to translate")

	}

	// Preparing the request data for the DeepL API
	www2URL := "https://www2.deepl.com/jsonrpc"
	id = id + 1
	postData := initDeepLXData(sourceLang, targetLang)
	text := Text{
		Text:                translateText,
		RequestAlternatives: 3,
	}
	postData.ID = id
	postData.Params.Texts = append(postData.Params.Texts, text)
	postData.Params.Timestamp = getTimeStamp(getICount(translateText))

	// Marshalling the request data to JSON and making necessary string replacements
	post_byte, _ := json.Marshal(postData)
	postStr := string(post_byte)

	// Adding spaces to the JSON string based on the ID to adhere to DeepL's request formatting rules
	if (id+5)%29 == 0 || (id+3)%13 == 0 {
		postStr = strings.Replace(postStr, "\"method\":\"", "\"method\" : \"", -1)
	} else {
		postStr = strings.Replace(postStr, "\"method\":\"", "\"method\": \"", -1)
	}

	// Creating a new HTTP POST request with the JSON data as the body
	post_byte = []byte(postStr)
	reader := bytes.NewReader(post_byte)
	request, err := http.NewRequest("POST", www2URL, reader)

	if err != nil {
		log.Println(err)
		return "", errors.New("Post request failed")
	}

	// Setting HTTP headers to mimic a request from the DeepL iOS App
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("x-app-os-name", "iOS")
	request.Header.Set("x-app-os-version", "16.3.0")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("x-app-device", "iPhone13,2")
	request.Header.Set("User-Agent", "DeepL-iOS/2.9.1 iOS 16.3.0 (iPhone13,2)")
	request.Header.Set("x-app-build", "510265")
	request.Header.Set("x-app-version", "2.9.1")
	request.Header.Set("Connection", "keep-alive")

	// Making the HTTP request to the DeepL API
	var client *http.Client
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return "", errors.New("Uknown error")
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", errors.New("DeepL API request failed")
	}
	defer resp.Body.Close()

	// Handling potential Brotli compressed response body
	var bodyReader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "br":
		bodyReader = brotli.NewReader(resp.Body)
	default:
		bodyReader = resp.Body
	}

	// Reading the response body and parsing it with gjson
	body, _ := io.ReadAll(bodyReader)
	// body, _ := io.ReadAll(resp.Body)
	var r req
	json.Unmarshal(body, &r)
	log.Printf("deeplx翻译返回的结构体:%+v\n", r)
	res := gjson.ParseBytes(body)
	log.Printf("%s\n", body)
	// Handling various response statuses and potential errors
	if res.Get("error.code").String() == "-32600" {
		log.Println(res.Get("error").String())
		return "", errors.New("Invalid target language")
	}

	var alternatives []string
	res.Get("result.texts.0.alternatives").ForEach(func(key, value gjson.Result) bool {
		alternatives = append(alternatives, value.Get("text").String())
		return true
	})
	if res.Get("result.texts.0.text").String() == "" {
		return "", errors.New("Translation failed, API returns an empty result.")
	} else {
		return res.Get("result.texts.0.text").String(), nil
	}
}
