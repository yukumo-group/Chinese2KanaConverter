# Chinese2KoeConverter

[![Go Test Workflow](https://github.com/yukumo-group/Chinese2KoeConverter/actions/workflows/test.yaml/badge.svg)](https://github.com/yukumo-group/Chinese2KoeConverter/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/badge/go-1.25.0-00ADD8.svg)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yukumo-group/Chinese2KanaConverter.svg)](https://pkg.go.dev/github.com/yukumo-group/Chinese2KanaConverter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

A Go library that converts **Chinese text into Japanese katakana (kana)**.

The converter treats Chinese characters phonetically: it reads them as pinyin and
transliterates that pinyin into kana. Non-Chinese parts of the text (English, existing
kana, punctuation, numbers, etc.) are transliterated separately so mixed-language
sentences survive the round trip.

```
你好        ->  ニーハオ
都会区      ->  トウホイチュイ
Hello!      ->  ヘロー！
```

## Features

- **Chinese → kana** by reading Chinese as pinyin and mapping each syllable to kana.
- **Two-stage pinyin-to-kana mapping**: a fast lookup table for complete pinyin syllables,
  with a fallback that splits a syllable into its *initial* and *final* and maps them
  independently so unseen syllables still convert.
- **Mixed-language handling**: text is split into Chinese / non-Chinese chunks. Chinese
  chunks go through pinyin, everything else is transliterated by
  [English2KanaTransliteration](https://github.com/Luigi-Pizzolito/English2KanaTransliteration).
- **Polyphonic (heteronym) words**: register custom pinyin for words that have an
  irregular pronunciation (e.g. `都会区`), either from a plain `map[string]string` or
  from a JSON file via the `Manager`.
- **Thread-safe** map loading and `Manager` access.

## How it works

```
Chinese text ─► SeparateToChunks ─┬─ IsChinese ─► ToPinyin ─► ToKana ─┐
                                  │                                    ├─► kana
                                  └─ not Chinese ─► OthersToKana ──────┘
```

1. `language.SeparateToChunks` walks the runes and groups them into runs that are either
   fully Chinese (matched by `^[\u4e00-\u9fa5\u3007]+$`) or not.
2. Chinese chunks are converted to pinyin with `cpyconverter.ToPinyin`. Tones are stripped
   because tones do not affect the kana output. A user-supplied heteronym dictionary is
   merged into `gpy`'s phrase dictionary first.
3. Each pinyin syllable is looked up in `table.PinyinToJapMap`. If it is missing, the
   syllable is matched against `table.InitialToKana` and `table.FinalToKana` (longest
   match first, see `table.SortStringSlice`).
4. Non-Chinese chunks are passed to `japconverter.OthersToKana`.


## Requirements

- [Go](https://go.dev/dl/) **1.25.0** or newer (see `go.mod`).

## Installation

The Go module path is `github.com/yukumo-group/Chinese2KanaConverter` (the repository is
named `Chinese2KoeConverter`).

```bash
go get github.com/yukumo-group/Chinese2KanaConverter
```

## Quick start

The main entry point is `converter.SingleChinesePieceToKana`.

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
)

func main() {
	// Set useHeteronym to true to honor registered polyphonic words.
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

The second parameter, `useHeteronym`, controls whether the custom polyphonic dictionary
is consulted when reading pinyin.

## Polyphonic (heteronym) words

Some words are read with a pronunciation that a per-character pinyin dictionary would
get wrong. You can override those readings.

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
	// Use SafeLoadPolyphonics when loading concurrently.
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

`polyphonic.Manager` keeps a dictionary of heteronyms and can read/write it as JSON.

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	manager := polyphonic.NewManager()

	// Add a word and make it available to the converter immediately.
	manager.AddPolyphonic("都会区", "du hui qu")

	// Persist to disk...
	manager.SetTargetFile("polyphonic.json")
	if err := manager.Save(); err != nil {
		panic(err)
	}

	// ...and load it back later.
	loaded, err := polyphonic.NewManagerFromFile("polyphonic.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(loaded.GetData()) // map[都会区:du hui qu]

	// Re-register everything stored in the manager with the converter.
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

> **Note:** `NewManager`, `NewManagerFromFile`, `GetData` and the other `Manager` methods
> are safe for concurrent use; the `Manager` embeds a `sync.RWMutex`.


## API reference

| Package | Symbol | Description |
| --- | --- | --- |
| `pkg/converter` | `SingleChinesePieceToKana(chineseText string, useHeteronym bool) (string, error)` | Converts a piece of text to kana. The main entry point. |
| `pkg/polyphonic` | `LoadPolyphonics(map[string]string)` | Registers a heteronym dictionary with the converter. |
| `pkg/polyphonic` | `SafeLoadPolyphonics(map[string]string)` | Same as above, guarded by a mutex for concurrent callers. |
| `pkg/polyphonic` | `NewManager() *Manager` | Creates an in-memory heteronym manager. |
| `pkg/polyphonic` | `NewManagerFromFile(path string) (*Manager, error)` | Loads a manager from a JSON file. |
| `pkg/polyphonic` | `(*Manager).AddPolyphonic(chinese, pinyin string)` | Adds or overwrites one heteronym. |
| `pkg/polyphonic` | `(*Manager).SetTargetFile(path string)` | Sets the file used by `Save`. |
| `pkg/polyphonic` | `(*Manager).Save() error` | Writes the manager to its target file as JSON. |
| `pkg/polyphonic` | `(*Manager).GetData() map[string]string` | Returns a copy of the heteronym map. |
| `pkg/polyphonic` | `(*Manager).Initialize()` | Registers every stored heteronym with the converter. |

Packages under `internal/` are implementation details and are not part of the public API.

## Project structure

```
.
├── go.mod                     # module github.com/yukumo-group/Chinese2KanaConverter
├── .goreleaser.yml            # GoReleaser config for tagged releases
├── internal/
│   ├── cpyconverter/          # Chinese -> pinyin (go-ego/gpy) + heteronym dictionary
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
git clone https://github.com/yukumo-group/Chinese2KoeConverter.git
cd Chinese2KoeConverter
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

