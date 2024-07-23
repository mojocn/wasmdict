# WASM ECDICT
A English-Chinese Oxford dictionary lib both for Go and WASM.


## Credits
Thanks for https://github.com/skywind3000/ECDICT 's dictionary data.


## How to build WASM

```bash
# build wasm
GOARCH=wasm GOOS=js go build -o ecdict.wasm
```

## How to use in JavaScript

```typescript
import "./wasm_exec.js"; # from  https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js

type DictLookupFn = (word: string) => WordItem | null;

declare global {
  export interface Window {
    Go: any;
    lookUp: DictLookupFn;
  }
}


async function loadWasm(wasmUrl: string = "/ecdict.wasm") {
  //https://davetayls.me/blog/2022-11-24-use-wasm-compiled-golang-functions-in-nextjs
  try {
    if ("lookUp" in window && typeof window.lookUp === "function") {
      return window.lookUp;
    }
    const go = new window.Go(); // Defined in wasm_exec.js
    let wasm: WebAssembly.WebAssemblyInstantiatedSource;
    if ("instantiateStreaming" in WebAssembly) {
      wasm = await WebAssembly.instantiateStreaming(
        fetch(wasmUrl),
        go.importObject,
      );
    } else {
      const resp = await fetch(wasmUrl);
      const bytes = await resp.arrayBuffer();
      wasm = await WebAssembly.instantiate(bytes, go.importObject);
    }
    go.run(wasm.instance);
    return window.lookUp;// now it's available, you can use `window.lookUp('Awesome')` to look up a word.
  } catch (e) {
    console.error(e);
    return null;
  }
}


```



## how to use in Go

`go get -u github.com/mojocn/wasmecdict`

```go
package wasmecdict

import (
	"testing"
)

func TestLookUp(t *testing.T) {
	for _, s := range []string{"awesome", "America", "Europe", "China", "book", "joker", "polish", "Polish", "china", "China"} {
		word := LookUp(s)
		if word == nil {
			t.Errorf("Word %s not found", s)
		} else {
			t.Logf("Word %s found: %s, %s", s, word.Definition, word.Translation)
		}
	}
}

```
