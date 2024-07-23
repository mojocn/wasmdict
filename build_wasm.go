//go:build js && wasm

package wasmecdict

import (
	"syscall/js"
)

const version = "0.0.1"

func lookUp(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		res := LookUp(text)
		if res != nil {
			return js.ValueOf(res.toMap())
		}
	}
	return js.Undefined()
}

func info(_ js.Value, _ []js.Value) any {
	return map[string]interface{}{
		"version":        version,
		"words":          len(dictMapSingleton),
		"lemmas":         len(lemmaMapSingleton),
		"author":         "EricZhou@mojotv.cn",
		"email":          "neochau@gmail.com",
		"DictionaryData": "https://github.com/skywind3000/ECDICT",
		"License":        "MIT",
	}
}

func main() {
	js.Global().Set("lookUp", js.FuncOf(lookUp))     //export window.lookUp to lookUp English word
	js.Global().Set("ecDictionary", js.FuncOf(info)) // export ecDictionary show info
	done := make(chan struct{}, 0)
	<-done
}
