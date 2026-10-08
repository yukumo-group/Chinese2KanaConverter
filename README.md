# Chinese2KanaConverter

**English** | [简体中文](./README.zh-CN.md) | [日本語](./README.ja-JP.md)

[![Go Test Workflow](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml/badge.svg)](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/badge/go-1.25.0-00ADD8.svg)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yukumo-group/Chinese2KanaConverter.svg)](https://pkg.go.dev/github.com/yukumo-group/Chinese2KanaConverter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

A Go library that converts **Chinese text into Japanese katakana (kana)**.

The converter reads Chinese characters phonetically: it turns them into pinyin and
transliterates that pinyin into kana. Non-Chinese parts of the text (English, existing
kana, punctuation, numbers, ...) are transliterated separately, so mixed-language
sentences are handled chunk by chunk instead of being thrown away.

```
你好                 ->  ニーハオ
都会区               ->  トウホイチュイ
Hello!               ->  ヘロー！
你好，世界！Hello!    ->  ニーハオシーチエ！ヘロー！
```

> **Note:** This is a *phonetic approximation* tuned for Japanese pronunciation, not a
> precise transcription of Chinese. Tones, aspiration and several syllable contrasts are
> intentionally collapsed into the nearest kana. See
> [Pronunciation model](#pronunciation-model) and [Limitations](#limitations) for details.

## Features

- **Chinese -> kana**: reads Chinese as pinyin and maps every syllable to kana.
- **Two-stage pinyin-to-kana mapping**: a lookup table for complete pinyin syllables,
  plus a fallback that splits unknown syllables into *initial* + *final* and maps each
  part, so every pinyin input still produces output.
- **Mixed-language aware**: the text is split into Chinese / non-Chinese chunks that are
  converted independently. Chinese chunks go through pinyin; everything else is
  transliterated by
  [English2KanaTransliteration](https://github.com/Luigi-Pizzolito/English2KanaTransliteration).
- **Polyphonic (heteronym) words**: register custom pinyin for words with an irregular
  pronunciation (e.g. `都会区`), either from a plain `map[string]string` or from a JSON
  file via the `Manager`.
- **Custom gse dictionaries**: load your own [gse](https://github.com/go-ego/gse)
  dictionary — or generate one from your polyphonic map — so the segmenter recognises the
  words that matter for your corpus.
- **Thread-safe** map loading and `Manager` access.

## How it works

```
Chinese text
     |
     v
SeparateToChunks ......... split into Chinese / non-Chinese runs
     |
     +-- Chinese ------> ToPinyin --> ToKana ------------------+
     |                   (gpy +      (table lookup, or        |
     |                    heteronym   initial + final)        |
     |                    dictionary)                          +--> kana
     |                                                         |
     +-- non-Chinese --> OthersToKana (English / kana) --------+
```

1. `language.SeparateToChunks` walks the runes and groups them into runs that are either
   fully Chinese (matched by `^[\u4e00-\u9fa5\u3007]+$`) or not.
2. Each Chinese chunk is converted to pinyin with `cpyconverter.ToPinyin`. Tones are
   stripped because they do not affect the kana output. Any user-supplied heteronym
   dictionary is merged into `gpy`'s phrase dictionary first, and an optional
   [custom gse dictionary](#custom-gse-dictionary) can steer how the text is segmented.
3. Each pinyin syllable is looked up in `table.PinyinToJapMap`. If it is missing, the
   syllable is matched against `table.InitialToKana` and `table.FinalToKana`
   (longest match first, see `table.SortStringSlice`).
4. Each non-Chinese chunk is passed to `japconverter.OthersToKana`.
5. The per-chunk results are concatenated in order.

## Requirements

- [Go](https://go.dev/dl/) **1.25.0** or newer (see `go.mod`).

## Installation

The Go module and the repository are both `github.com/yukumo-group/Chinese2KanaConverter`.

```bash
go get github.com/yukumo-group/Chinese2KanaConverter
```

## Quick start

The entry point is `converter.SingleChinesePieceToKana`.

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
)

func main() {
	// The second argument enables the custom polyphonic dictionary.
	kana, err := converter.SingleChinesePieceToKana("你好", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(kana) // ニーハオ
}
```

Output:

```
ニーハオ
```

The second parameter, `useHeteronym`, controls whether the registered polyphonic
dictionary is consulted when reading pinyin.

## Mixed-language input

Text is converted chunk by chunk, so Chinese and non-Chinese parts can be freely mixed:

```go
kana, err := converter.SingleChinesePieceToKana("你好，世界！Hello!", true)
// kana == "ニーハオシーチエ！ヘロー！"
```

Every chunk is converted on its own and the results are concatenated:

| Chunk | Type | Output |
| --- | --- | --- |
| `你好` | Chinese | `ニーハオ` |
| `，` | non-Chinese | *(dropped by the transliterator)* |
| `世界` | Chinese | `シーチエ` |
| `！` | non-Chinese | `！` |
| `Hello` | non-Chinese | `ヘロー` |
| `!` | non-Chinese | `！` |

## Polyphonic (heteronym) words

Some words are read with a pronunciation a per-character dictionary would get wrong.
Register those readings so the converter picks them up.

### Load a static map directly

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	// Chinese word -> pinyin (syllables separated by spaces).
	// Use SafeLoadPolyphonics when loading from several goroutines.
	polyphonic.LoadPolyphonics(map[string]string{
		"都会区": "du hui qu",
	})

	kana, err := converter.SingleChinesePieceToKana("都会区", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(kana) // トウホイチュイ
}
```

### Persist polyphonics with the Manager

`polyphonic.Manager` holds a heteronym dictionary and can read/write it as JSON.

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	manager := polyphonic.NewManager()

	// Add a word; it becomes available to the converter immediately.
	manager.AddPolyphonic("都会区", "du hui qu")

	// The manager persists two files: the heteronyms as JSON, and a generated gse
	// dictionary used to segment the text.
	manager.SetTargetFile("polyphonic.json", "polyphonic.dict.txt")
	if err := manager.Save(); err != nil {
		panic(err)
	}
	if err := manager.SaveGSEDict(); err != nil {
		panic(err)
	}

	// ...and load both back later. The gse dictionary is registered automatically.
	loaded, err := polyphonic.NewManagerFromFile("polyphonic.json", "polyphonic.dict.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(loaded.GetData()) // map[都会区:du hui qu]

	// Re-register everything the manager holds with the converter.
	loaded.Initialize()
}
```

The JSON file written by `Save` looks like this:

```json
{
  "heteronym": {
    "都会区": "du hui qu"
  }
}
```

> **Note:** `Manager` embeds a `sync.RWMutex`, so `NewManager`, `NewManagerFromFile`,
> `GetData` and the other methods are safe for concurrent use.

## Custom gse dictionary

Phrase-level pinyin is resolved by `go-ego/gpy` *after* the text has been segmented by
[gse](https://github.com/go-ego/gse). Every line of a gse dictionary is an entry
`word frequency [part-of-speech]`, for example:

```
都会区 100 n
```

Such a dictionary only influences **segmentation**: the pinyin of a segmented word still has
to be known to `gpy`, either from its built-in phrase dictionary or from a registered
[heteronym](#polyphonic-heteronym-words). Otherwise the word falls back to per-character
readings.

### Loading a gse dictionary

```go
import "github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"

// Accepts one or more gse dictionary files.
if err := converter.InitGSEDict("my-dict.txt"); err != nil {
	panic(err)
}
```

`InitGSEDict` accepts several files at once, so a base dictionary and your own can be
loaded in a single call.

### Writing a gse dictionary

A gse dictionary can be generated from a polyphonic map, so the same words you registered
are also the words the segmenter recognises:

```go
import "github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"

polyphonics := map[string]string{
	"都会区": "du hui qu",
	"西雅图": "xi ya tu",
}

// Write a standalone gse dictionary file.
if err := polyphonic.WriteDict(polyphonics, "polyphonic.dict.txt"); err != nil {
	panic(err)
}
```

`WriteDict` writes only the **Chinese words** (the pinyin stays in the JSON heteronym file),
with a fixed frequency of `100` and the part-of-speech `n`. For the map above the generated
file is:

```
都会区 100 n
西雅图 100 n
```

When everything is driven through a `Manager`, call `SaveGSEDict` instead. It writes the
dictionary to the path previously handed to `SetTargetFile`:

```go
manager.SetTargetFile("polyphonic.json", "polyphonic.dict.txt")

if err := manager.Save(); err != nil { // heteronyms -> JSON
	panic(err)
}
if err := manager.SaveGSEDict(); err != nil { // words -> gse dictionary
	panic(err)
}
```

> **Note:** `NewManager` already defaults to `polyphonic.json` and `dict.txt`, so
> `SetTargetFile` is only needed when you want other paths. `NewManagerFromFile` creates the
> JSON file and the gse dictionary when they do not exist yet, and registers the dictionary
> with the converter for you.

## Pronunciation model

The output is an *approximation* of Chinese using Japanese phonology, not an exact
transcription. The following simplifications are intentional and come straight from the
mapping tables in `internal/table`:

| Aspect | Behaviour | Example |
| --- | --- | --- |
| Tones | dropped (`gpy.Normal`) | 妈 / 麻 / 马 / 骂 -> `マー` |
| Aspiration | unaspirated and aspirated stops share kana | `ba` = `pa` = `パー` |
| Retroflex vs. palatal | many initials collapse | `zhi` / `chi` / `ji` / `qi` -> `チー`; `shi` / `xi` -> `シー` |
| Nasal finals | `-n` / `-ng` both become `ン` | `yin` = `ying` = `イン` |
| Vowels | Chinese `e` maps like `o` | `e` -> `オー` |
| Syllable length | single syllables often get a long-vowel mark | `ta` -> `ター` |

Because of this, different characters can produce the same kana. If you need stricter
readings, tune the tables in `internal/table/jap_to_chn.go` and/or register heteronyms.

## Limitations

The library targets **Chinese -> kana**. English and Japanese *kana* work as a side
effect of the chunking, but **Japanese kanji is not handled correctly**.

| Input language | Supported | Example |
| --- | --- | --- |
| Chinese | Yes, approximate (see [Pronunciation model](#pronunciation-model)) | `我爱你` -> `ウォーアイニー` |
| English | Yes | `Hello` -> `ヘロー` |
| Japanese kana (hiragana / katakana) | Yes | `こんにちは` -> `コンニチハ` |
| Japanese kanji | **No** | `日本語` -> `リーペンユイ`, `東京` -> `トンチン` |

### Why Japanese kanji is wrong

`language.IsChinese` classifies a run as Chinese when every rune matches
`^[\u4e00-\u9fa5\u3007]+$`. Japanese kanji sit in the **same Unicode block** as Chinese
characters (U+4E00–U+9FA5), so at the code-point level they cannot be told apart and are
routed to the pinyin path (`cpyconverter.ToPinyin`). The output is therefore the
*Chinese* reading, not the Japanese one:

- `日本語` -> `リーペンユイ` (rì běn yǔ) instead of `ニホンゴ`
- `東京` -> `トンチン` (dōng jīng) instead of `トウキョウ`

The dependency `English2KanaTransliteration` actually ships a Japanese kanji reader
(`Kanji_to_kana.go`, backed by `dict/kanjidic2_pronounce.json`) and it is wired up through
`japconverter.OthersToKana`. A pure-kanji run never reaches it because of the routing
above. Note that even when it is reached, its readings are looked up **per character
without context**, so `日本` would become `ニチホン` rather than `ニホン`. For reference,
calling `OthersToKana` directly gives:

| Input | `language.IsChinese` | `OthersToKana` |
| --- | --- | --- |
| `日本語` | `true` | `ニチホンゴ` |
| `東京` | `true` | `トウキョウ` |

### Workarounds

- **Decide the language before calling**: split the text by language yourself and feed
  only the Chinese parts to this library. Convert Japanese segments with a dedicated
  reading/furigana tool first, then pass the resulting kana in (kana is passed through).
- **Pre-convert kanji to kana**: since kana runs are preserved, converting Japanese kanji
  to kana (furigana) before conversion sidesteps the misclassification entirely.
- Registering the word as a heteronym does **not** help here: heteronym values are pinyin
  and still go through the pinyin path, so they cannot produce an arbitrary Japanese
  reading.

## API reference

| Package | Symbol | Description |
| --- | --- | --- |
| `pkg/converter` | `SingleChinesePieceToKana(chineseText string, useHeteronym bool) (string, error)` | Converts a piece of text to kana. The main entry point. |
| `pkg/converter` | `InitGSEDict(paths ...string) error` | Loads one or more custom gse dictionary files used for word segmentation. |
| `pkg/polyphonic` | `LoadPolyphonics(map[string]string)` | Registers a heteronym dictionary with the converter. |
| `pkg/polyphonic` | `SafeLoadPolyphonics(map[string]string)` | Same as above, guarded by a mutex for concurrent callers. |
| `pkg/polyphonic` | `NewManager() *Manager` | Creates an in-memory heteronym manager. Defaults to `polyphonic.json` and `dict.txt`. |
| `pkg/polyphonic` | `NewManagerFromFile(targetFilePath, targetDictPath string) (*Manager, error)` | Loads a manager from a JSON file and registers the gse dictionary. |
| `pkg/polyphonic` | `(*Manager).AddPolyphonic(chinese, pinyin string)` | Adds or overwrites one heteronym. |
| `pkg/polyphonic` | `(*Manager).SetTargetFile(targetFilePath, targetDictPath string)` | Sets the JSON file and the gse dictionary used by `Save` / `SaveGSEDict`. |
| `pkg/polyphonic` | `(*Manager).Save() error` | Writes the manager to its target file as JSON (default `polyphonic.json`). |
| `pkg/polyphonic` | `(*Manager).SaveGSEDict() error` | Writes the heteronyms out as a gse dictionary file (default `dict.txt`). |
| `pkg/polyphonic` | `WriteDict(polyphonics map[string]string, targetFilePath string) error` | Writes a gse dictionary file from a Chinese -> pinyin map. |
| `pkg/polyphonic` | `(*Manager).GetData() map[string]string` | Returns a copy of the heteronym map. |
| `pkg/polyphonic` | `(*Manager).Initialize()` | Registers every stored heteronym with the converter. |

Packages under `internal/` are implementation details and are not part of the public API.

## Project structure

```
.
├── go.mod                     # module github.com/yukumo-group/Chinese2KanaConverter
├── .goreleaser.yml            # GoReleaser config for tagged releases
├── internal/
│   ├── cpyconverter/          # Chinese -> pinyin (go-ego/gpy), heteronym + gse dicts
│   ├── japconverter/          # pinyin -> kana, and non-Chinese -> kana
│   ├── language/              # text chunking and Chinese detection
│   ├── process/               # small shared helpers
│   └── table/                 # pinyin / initial / final -> kana mapping tables
└── pkg/
    ├── converter/             # public entry point
    └── polyphonic/            # heteronym Manager and loaders
```

## Development

Clone the repository and run the test suite:

```bash
git clone https://github.com/yukumo-group/Chinese2KanaConverter.git
cd Chinese2KanaConverter
go mod download
go test ./...
```

Run tests the same way CI does (verbose, with coverage):

```bash
go test -v ./... -tags=test -coverpkg=./pkg/...,./internal/... -coverprofile=coverage.txt
```

Static analysis and lint (as used in CI):

```bash
go vet ./...
go install golang.org/x/lint/golint@latest
golint -set_exit_status ./...
```

## Continuous integration

- [`.github/workflows/test.yaml`](./.github/workflows/test.yaml) runs tests, `go vet` and
  `golint` on every push and pull request, and uploads coverage to Codecov.
- [`.github/workflows/release.yaml`](./.github/workflows/release.yaml) runs GoReleaser
  when a tag matching `v*` is pushed.

## Contributing

1. Fork the repository and create a feature branch.
2. Keep the existing style: tabs for indentation, doc comments on exported symbols
   (golint must pass), and table-driven tests in `*_test.go`.
3. Add or update tests for any behavior change.
4. Make sure `go test ./...`, `go vet ./...` and `golint -set_exit_status ./...` all pass.
5. Open a pull request.

## License

Released under the [MIT License](./LICENSE). Copyright (c) 2026 Yukumo Group.

