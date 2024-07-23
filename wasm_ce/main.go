//go:build js && wasm

package main

import (
	"github.com/mojocn/wasmdict"
	"syscall/js"
)

const versionCe = "1.0.1"

func ceLookUpTw(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		res := wasmdict.CeLookUp(text, false)
		if res != nil {
			return js.ValueOf(res.Map())
		}
	}
	return js.Undefined()
}
func ceLookUpCn(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		res := wasmdict.CeLookUp(text, true)
		if res != nil {
			return js.ValueOf(res.Map())
		}
	}
	return js.Undefined()
}
func ceQueryLikeTw(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		words := wasmdict.CeQueryLike(text, false, 10)
		results := make([]map[string]interface{}, 0, len(words))
		for _, word := range words {
			results = append(results, word.Map())
		}
		return js.ValueOf(results)
	}
	return js.Undefined()
}
func ceQueryLikeCn(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		words := wasmdict.CeQueryLike(text, true, 10)
		results := make([]map[string]interface{}, 0, len(words))
		for _, word := range words {
			results = append(results, word.Map())
		}
		return js.ValueOf(results)
	}
	return js.Undefined()
}

func ceDictInfo(_ js.Value, _ []js.Value) any {
	return map[string]interface{}{
		"version":        versionCe,
		"author":         "EricZhou@mojotv.cn",
		"email":          "neochau@gmail.com",
		"DictionaryData": "https://www.mdbg.net/chinese/dictionary?page=cedict",
		"License":        "MIT",
	}
}

func main() {
	wasmdict.PreLoadCeDict()
	js.Global().Set("ceLookUpTw", js.FuncOf(ceLookUpTw))       //export window.ceLookUpTw to lookUp Traditional Chinese word
	js.Global().Set("ceLookUpCn", js.FuncOf(ceLookUpCn))       //export window.ceLookUpCn to lookUp Simplified Chinese word
	js.Global().Set("ceQueryLikeCn", js.FuncOf(ceQueryLikeCn)) //export window.ceQueryLikeCn to query Simplified Chinese word
	js.Global().Set("ceQueryLikeTw", js.FuncOf(ceQueryLikeTw)) //export window.ceQueryLikeTw to query Traditional Chinese word
	js.Global().Set("ceDictInfo", js.FuncOf(ceDictInfo))       //export window.ceDictInfo to get dictionary info
	done := make(chan struct{})
	<-done
}
