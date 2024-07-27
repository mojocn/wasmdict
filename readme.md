[![Go Reference](https://pkg.go.dev/badge/github.com/mojocn/wasmdict.svg)](https://pkg.go.dev/github.com/mojocn/wasmdict)




# WASM ECDICT
A English-Chinese Oxford dictionary lib both for `Go` and `WebAssembley`.

**[LiveDemo English-Chinese Dictionary](https://mojotv.cn/gadgets/english-chinese-dictionary)**

## Credits
Thanks for https://github.com/skywind3000/ECDICT 's dictionary data.


## How to build WASM

```bash
# build wasm English-Chinese dictionary
GOARCH=wasm GOOS=js go build -o dict_ec.wasm wasm_ec/main.go

# build wasm Chinese-English dictionary
GOARCH=wasm GOOS=js go build -o dict_ce.wasm wasm_ce/main.go
```

or just use the prebuilt wasm file from https://github.com/mojocn/wasmecdict/releases/download/v1.0.1/ecdict.wasm

## How to use in JavaScript

### English-Chinese Dictionary WASM Usage

```typescript
import "./wasm_exec.js"; # from  https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js

export interface WordEntry {
    word: string;
    phonetic: string;
    definition: string;
    translation: string;
    pos: string;
    collins: string;
    oxford: string;
    tag: string;
    bnc: string;
    frq: string;
    exchange: string;
    detail: string;
    audio: string;
}


declare global {
    export interface Window extends EcDict {
        Go: any;
    }
}

interface EcDict {
    ecLookUp: (word: string) => WordEntry | undefined;
    ecQueryLike: (word: string) => string[] | undefined;
    ecDictInfo: () => Object;
}

async function loadWasmEc(wasmUrl: string = "/dict_ec.wasm") {
    try {
        if ("ecLookUp" in window && typeof window.ecLookUp === "function") {
            return {
                ecLookUp: window.ecLookUp,
                ecQueryLike: window.ecQueryLike,
                ecDictInfo: window.ecDictInfo
            } as EcDict;
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
        return {
            ecLookUp: window.ecLookUp,
            ecQueryLike: window.ecQueryLike,
            ecDictInfo: window.ecDictInfo
        } as EcDict;
    } catch (e) {
        console.error(e);
        return null;
    }
}


```

usage:

```typescript

const wordInfo = window.ecLookUp('Awesome');
console.log(wordInfo);

```

### Chinese-English Dictionary WASM Usage

```typescript
import "./wasm_exec.js"; # from  https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js

export interface WordEntryCe {
    traditional: string;
    simplified: string;
    pinyin: string;
    english: string;
}


declare global {
    export interface Window extends DictCe {
        Go: any;
    }
}

interface DictCe {
    ceLookUp: (word: string, isZhCn: boolean) => WordEntryCe | undefined;
    ceQueryLike: (zi: string, isZhCn: boolean) => WordEntryCe[];
    ceDictInfo: () => Object;
}


async function loadWasmCe(wasmUrl: string = "/dict_ce.wasm") {
    try {
        if ("ceLookUp" in window && typeof window.ceLookUp === "function") {
            return {
                ceLookUp: window.ceLookUp,
                ceQueryLike: window.ceQueryLike,
                ceDictInfo: window.ceDictInfo
            } as DictCe;
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
        return {
            ceLookUp: window.ceLookUp,
            ceQueryLike: window.ceQueryLike,
            ceDictInfo: window.ceDictInfo
        } as DictCe;
    } catch (e) {
        console.error(e);
        return null;
    }
}


```



## how to use in Go

`go get -u github.com/mojocn/wasmecdict`


```
