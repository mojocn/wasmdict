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

`go get -u github.com/mojocn/wasmdict`

```


"{"role": "assistant", "type": "code", "content": "{"function": "websearch", "query": "2024年 奥运会 争议事件", "prompt": "根据\\\\"2024年 奥运会 争议事件\\\\"的搜索结果,提取有关2024年奥运会的争议事件信息", "query_tag": ["体育", "新闻", "国际"]}", "file_ids": null, "show_content": "[object Object],[object Object],[object Object],[object Object],[object Object],[object Object],[object Object],[object Object],[object Object]**美国队服设计争议**：2024年巴黎奥运会开幕式上，美国队的队服设计引起了广泛争议。拉尔夫·劳伦设计的队服包括海军蓝羊绒大衣和红白相间的条纹，以及小马图案和美国奥运代表队徽章，但这一设计并未获得预期好评，反而在社会各阶层及网络舆论中引起广泛争议。高昂的售价（高达1000美元）也引来了消费者的异议与争议，最终被誉为“史上最差设计”[^1]。\\n\\n**奥运村大巴空调问题**：巴黎奥运村的运动员大巴车没有空调，导致一名韩国选手在下车后直接晕倒。这一事件引起了公众的广泛关注和批评，认为奥运村未为运动员提供配备空调的大巴是不可接受的。巴黎奥运会紧急调派了带有空调设施的客车替代原先的大巴[^2]。\\n\\n**国旗失误争议**：在2024年奥运会上，韩国国旗在直播中被虚化，引发了韩国网友的愤怒和不满。这一事件反映了国际赛事组织中的技术失误和文化敏感度缺失，以及对国旗象征性元素的错误处理[^3]。\\n\\n**国际奥委会制裁争议**：国际奥委会对俄罗斯和白俄罗斯运动员实施制裁，禁止他们在2024年奥运会开幕式上列队游行。这一决定被指责为“种族主义和新纳粹主义”，引发了国际社会的广泛争议[^4]。\\n\\n**耐克美国田径队服争议**：耐克发布的2024年奥运会美国田径队官方队服，尤其是女子服装设计，引发了批评。有观点认为，女子服装设计过于暴露，不符合运动员的舒适和专业形象[^5]。\\n\\n这些争议事件反映了奥运会组织和管理中存在的一些问题，包括对文化差异的敏感度、技术失误、以及对运动员权益的考虑等方面。这些问题需要国际奥委会和相关组织认真对待和改进，以确保奥运会的顺利进行和运动员的权益得到保障。"}"
