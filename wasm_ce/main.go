//go:build js && wasm

package main

import (
	"github.com/mojocn/wasmdict"
	"syscall/js"
)

const versionCe = "1.0.1"

func ceLookUp(_ js.Value, args []js.Value) any {
	if len(args) > 1 {
		isCnZh := args[1].Bool()
		text := args[0].String()
		res := wasmdict.CeLookUp(text, isCnZh)
		if res != nil {
			return js.ValueOf(res.Map())
		}
	}
	return js.Undefined()
}
func ceQueryLike(_ js.Value, args []js.Value) any {
	if len(args) > 1 {
		text := args[0].String()
		isCnZh := args[1].Bool()
		words := wasmdict.CeQueryLike(text, isCnZh, 10)
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
	js.Global().Set("ceLookUp", js.FuncOf(ceLookUp))
	js.Global().Set("ceQueryLike", js.FuncOf(ceQueryLike)) //
	js.Global().Set("ceDictInfo", js.FuncOf(ceDictInfo))   //export window.ceDictInfo to get dictionary info
	done := make(chan struct{})
	<-done
}
