//go:build js && wasm

package main

import (
	"github.com/mojocn/wasmdict"
	"syscall/js"
)

const versionEc = "1.0.1"

func lookUp(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		res := wasmdict.EcLookUp(text)
		if res != nil {
			return js.ValueOf(res.Map())
		}
	}
	return js.Undefined()
}

func queryLike(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		words := wasmdict.EcQueryLike(text, 10)
		results := make([]interface{}, len(words))
		for i, word := range words {
			results[i] = word
		}
		return js.ValueOf(results)
	}
	return js.Undefined()

}

func ecDictInfo(_ js.Value, _ []js.Value) any {
	return map[string]interface{}{
		"version":        versionEc,
		"author":         "EricZhou@mojotv.cn",
		"email":          "neochau@gmail.com",
		"DictionaryData": "https://github.com/skywind3000/ECDICT",
		"License":        "MIT",
	}
}

func main() {
	wasmdict.PreLoadEcDict()
	js.Global().Set("ecQueryLike", js.FuncOf(queryLike)) //use window.ecQueryLike("word") to query like a word
	js.Global().Set("ecLookUp", js.FuncOf(lookUp))       //use window.ecLookUp("word") to look up a word
	js.Global().Set("ecDictInfo", js.FuncOf(ecDictInfo)) //use window.ecDictInfo() to get dictionary info
	done := make(chan struct{})
	<-done
}
