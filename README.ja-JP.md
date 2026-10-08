# Chinese2KanaConverter

[English](./README.md) | [简体中文](./README.zh-CN.md) | **日本語**

[![Go Test Workflow](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml/badge.svg)](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/badge/go-1.25.0-00ADD8.svg)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yukumo-group/Chinese2KanaConverter.svg)](https://pkg.go.dev/github.com/yukumo-group/Chinese2KanaConverter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

**中国語のテキストを日本語のカタカナ（かな）に変換する** Go ライブラリです。

本コンバーターは中国語の漢字を**音声的**に扱います。まず漢字をピンインに変換し、そのピンインを
かなに音訳します。テキスト中の非中国語部分（英語、既存のかな、記号、数字など）は個別に音訳
されるため、多言語が混在した文も**チャンク単位**で処理され、捨てられません。

```
你好                 ->  ニーハオ
都会区               ->  トウホイチュイ
Hello!               ->  ヘロー！
你好，世界！Hello!    ->  ニーハオシーチエ！ヘロー！
```

> **注意：** これは日本語の音韻に合わせた*近似的な音訳*であり、中国語の正確な転写では
> ありません。声調、有気音、いくつかの音節の対立は、意図的に最も近いかなへ統合されます。
> 詳細は[発音モデル](#発音モデル)と[制限事項](#制限事項)を参照してください。

## 機能

- **中国語 -> かな**：中国語をピンインとして読み、各音節をかなにマッピングします。
- **2 段階のピンイン→かなマッピング**：完全なピンイン音節用のルックアップテーブルに加え、
  未知の音節を*声母* + *韻母*に分割して個別にマッピングするフォールバックを備えているため、
  どのピンイン入力でも出力が得られます。
- **多言語を考慮したチャンク処理**：テキストは中国語 / 非中国語のチャンクに分割され、
  それぞれ独立に変換されます。中国語チャンクはピンインを経由し、それ以外は
  [English2KanaTransliteration](https://github.com/Luigi-Pizzolito/English2KanaTransliteration)
  によって音訳されます。
- **多音字（異読語）**：不規則な読みを持つ語（例：`都会区`）に独自のピンインを登録できます。
  素の `map[string]string` を渡す方法と、`Manager` 経由で JSON ファイルを使う方法があります。
- **独自の gse 辞書**：自分専用の [gse](https://github.com/go-ego/gse) 辞書を読み込むか、
  多音字マップから生成できます。これによりセグメンタが自分のコーパスに固有の語を認識します。
- **スレッドセーフ**な辞書の読み込みと `Manager` へのアクセス。

## 仕組み

```
中国語テキスト
     |
     v
SeparateToChunks ......... 中国語 / 非中国語の連続部分に分割
     |
     +-- 中国語 ------> ToPinyin --> ToKana ------------------+
     |                  (gpy +      (テーブル参照、または      |
     |                   多音字辞書)  声母 + 韻母)             |
     |                                                         +--> かな
     |                                                         |
     +-- 非中国語 ----> OthersToKana（英語 / かな）-----------+
```

1. `language.SeparateToChunks` は文字単位で走査し、`^[\u4e00-\u9fa5\u3007]+$` に完全一致する
   （＝すべて中国語の）連続部分、または中国語を 1 文字も含まない連続部分をチャンクにまとめます。
2. 各中国語チャンクは `cpyconverter.ToPinyin` でピンインに変換されます。声調はかな出力に
   影響しないため除去されます。ユーザーが登録した多音字辞書はまず `gpy` のフレーズ辞書に
   マージされ、任意で[独自の gse 辞書](#独自の-gse-辞書)によって分かち書きを制御できます。
3. 各ピンイン音節は `table.PinyinToJapMap` で参照されます。見つからない場合は
   `table.InitialToKana` と `table.FinalToKana` に対してマッチングされます
   （長い韻母が優先。`table.SortStringSlice` を参照）。
4. 各非中国語チャンクは `japconverter.OthersToKana` に渡されます。
5. 各チャンクの結果は元の順序で連結されます。

## 動作要件

- [Go](https://go.dev/dl/) **1.25.0** 以上（`go.mod` を参照）。

## インストール

Go モジュールとリポジトリはどちらも `github.com/yukumo-group/Chinese2KanaConverter` です。

```bash
go get github.com/yukumo-group/Chinese2KanaConverter
```

## クイックスタート

エントリポイントは `converter.SingleChinesePieceToKana` です。

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
)

func main() {
	// 第 2 引数で独自の多音字辞書を有効にします。
	kana, err := converter.SingleChinesePieceToKana("你好", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(kana) // ニーハオ
}
```

出力：

```
ニーハオ
```

第 2 引数 `useHeteronym` は、ピンインを読むときに登録済みの多音字辞書を参照するかどうかを
制御します。

## 多言語が混在した入力

テキストはチャンク単位で変換されるため、中国語と非中国語を自由に混在させられます：

```go
kana, err := converter.SingleChinesePieceToKana("你好，世界！Hello!", true)
// kana == "ニーハオシーチエ！ヘロー！"
```

各チャンクは個別に変換され、結果が順に連結されます：

| チャンク | 種類 | 出力 |
| --- | --- | --- |
| `你好` | 中国語 | `ニーハオ` |
| `，` | 非中国語 | *（音訳器によって破棄）* |
| `世界` | 中国語 | `シーチエ` |
| `！` | 非中国語 | `！` |
| `Hello` | 非中国語 | `ヘロー` |
| `!` | 非中国語 | `！` |

## 多音字（異読語）

文字単位の辞書では正しく読めない語があります。そのような読みを登録すると、コンバーターが
それを参照するようになります。

### 静的マップを直接読み込む

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	// 中国語の語 -> ピンイン（音節はスペース区切り）。
	// 複数の goroutine から読み込む場合は SafeLoadPolyphonics を使ってください。
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

### Manager で多音字を保存する

`polyphonic.Manager` は多音字辞書を保持し、JSON として読み書きできます。

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	manager := polyphonic.NewManager()

	// 語を追加すると、すぐにコンバーターから参照できるようになります。
	manager.AddPolyphonic("都会区", "du hui qu")

	// Manager は 2 つのファイルを永続化します：多音字の JSON と、
	// 分かち書きに使う生成済みの gse 辞書です。
	manager.SetTargetFile("polyphonic.json", "polyphonic.dict.txt")
	if err := manager.Save(); err != nil {
		panic(err)
	}
	if err := manager.SaveGSEDict(); err != nil {
		panic(err)
	}

	// ……後で両方を読み戻します。gse 辞書は自動的に登録されます。
	loaded, err := polyphonic.NewManagerFromFile("polyphonic.json", "polyphonic.dict.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(loaded.GetData()) // map[都会区:du hui qu]

	// Manager が保持する内容をコンバーターに再登録します。
	loaded.Initialize()
}
```

`Save` が書き出す JSON は次のとおりです：

```json
{
  "heteronym": {
    "都会区": "du hui qu"
  }
}
```

> **注意：** `Manager` は `sync.RWMutex` を埋め込んでいるため、`NewManager`、
> `NewManagerFromFile`、`GetData` などのメソッドは並行に呼び出しても安全です。

## 独自の gse 辞書

フレーズレベルのピンインは、テキストが [gse](https://github.com/go-ego/gse) によって
**分かち書き**された*後に*、`go-ego/gpy` が決定します。gse 辞書の各行は
`word frequency [part-of-speech]` というエントリです。例：

```
都会区 100 n
```

この辞書は**分かち書き**にのみ影響します。分割された語のピンインは、やはり `gpy` が知って
いる必要があります（組み込みのフレーズ辞書、または登録済みの
[多音字](#多音字異読語)）。そうでない場合、その語は 1 文字ずつの読みにフォールバックします。

### gse 辞書を読み込む

```go
import "github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"

// 1 つ以上の gse 辞書ファイルを指定できます。
if err := converter.InitGSEDict("my-dict.txt"); err != nil {
	panic(err)
}
```

`InitGSEDict` は複数のファイルを一度に受け取れるので、基本辞書と自分の辞書を 1 回の呼び出しで
読み込めます。

### gse 辞書を書き出す

gse 辞書は多音字マップから生成できます。こうすると、登録した語がそのままセグメンタの認識する
語になります：

```go
import "github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"

polyphonics := map[string]string{
	"都会区": "du hui qu",
	"西雅图": "xi ya tu",
}

// 単体の gse 辞書ファイルを書き出します。
if err := polyphonic.WriteDict(polyphonics, "polyphonic.dict.txt"); err != nil {
	panic(err)
}
```

`WriteDict` は**中国語の語のみ**を書き出します（ピンインは JSON の多音字ファイルに残ります）。
頻度は `100`、品詞は `n` に固定されています。上記のマップから生成されるファイルは：

```
都会区 100 n
西雅图 100 n
```

すべてを `Manager` 経由で扱う場合は、代わりに `SaveGSEDict` を呼びます。これは
`SetTargetFile` に渡したパスへ辞書を書き出します：

```go
manager.SetTargetFile("polyphonic.json", "polyphonic.dict.txt")

if err := manager.Save(); err != nil { // 多音字 -> JSON
	panic(err)
}
if err := manager.SaveGSEDict(); err != nil { // 語 -> gse 辞書
	panic(err)
}
```

> **注意：** `NewManager` は既定で `polyphonic.json` と `dict.txt` を使うため、
> `SetTargetFile` は別のパスを使いたいときだけ必要です。`NewManagerFromFile` は JSON ファイルと
> gse 辞書が存在しない場合に自動生成し、辞書をコンバーターに登録します。

## 発音モデル

出力は日本語の音韻を用いた中国語の*近似*であり、正確な転写ではありません。以下の単純化は
意図的なもので、`internal/table` のマッピングテーブルに由来します：

| 項目 | 挙動 | 例 |
| --- | --- | --- |
| 声調 | 破棄（`gpy.Normal`） | 妈 / 麻 / 马 / 骂 -> `マー` |
| 有気 / 無気 | 有気音と無気音が同じかなを共有 | `ba` = `pa` = `パー` |
| そり舌音 vs 硬口蓋音 | 多くの声母が統合される | `zhi` / `chi` / `ji` / `qi` -> `チー`；`shi` / `xi` -> `シー` |
| 鼻音韻尾 | `-n` / `-ng` がどちらも `ン` になる | `yin` = `ying` = `イン` |
| 母音 | 中国語の `e` は `o` と同様に扱われる | `e` -> `オー` |
| 音節の長さ | 単音節には長音符が付くことが多い | `ta` -> `ター` |

このため、異なる漢字が同じかなになることがあります。より厳密な読みが必要な場合は、
`internal/table/jap_to_chn.go` のテーブルを調整するか、多音字を登録してください。

## 制限事項

本ライブラリは**中国語 -> かな**を対象としています。英語と日本語の*かな*はチャンク処理の
副産物として動作しますが、**日本語の漢字は正しく扱えません**。

| 入力言語 | 対応 | 例 |
| --- | --- | --- |
| 中国語 | 対応（近似。詳細は[発音モデル](#発音モデル)） | `我爱你` -> `ウォーアイニー` |
| 英語 | 対応 | `Hello` -> `ヘロー` |
| 日本語のかな（ひらがな / カタカナ） | 対応 | `こんにちは` -> `コンニチハ` |
| 日本語の漢字 | **非対応** | `日本語` -> `リーペンユイ`、`東京` -> `トンチン` |

### 日本語の漢字が誤って読まれる理由

`language.IsChinese` は「すべての文字が `^[\u4e00-\u9fa5\u3007]+$` に一致する」ときにその
連続部分を中国語と判定します。日本語の漢字は中国語の漢字と**同じ Unicode ブロック**
（U+4E00–U+9FA5）にあるため、コードポイントのレベルでは区別できず、ピンインの経路
（`cpyconverter.ToPinyin`）へルーティングされます。その結果、日本語の読みではなく**中国語の
読み**になります：

- `日本語` -> `リーペンユイ`（rì běn yǔ）。正しくは `ニホンゴ`
- `東京` -> `トンチン`（dōng jīng）。正しくは `トウキョウ`

依存ライブラリの `English2KanaTransliteration` は、実は日本語漢字のリーダー
（`Kanji_to_kana.go`、`dict/kanjidic2_pronounce.json` を使用）を同梱しており、
`japconverter.OthersToKana` 経由で組み込まれています。しかし漢字だけのチャンクは上記の
ルーティングのため**そこへ到達しません**。また、仮に到達しても読みは**文脈を考慮しない
1 文字単位のテーブル参照**であるため、`日本` は `ニホン` ではなく `ニチホン` になります。
参考として、`OthersToKana` を直接呼び出した結果は次のとおりです：

| 入力 | `language.IsChinese` | `OthersToKana` |
| --- | --- | --- |
| `日本語` | `true` | `ニチホンゴ` |
| `東京` | `true` | `トウキョウ` |

### 回避策

- **呼び出し前に言語を決める**：テキストを言語ごとに自分で分割し、中国語の部分だけを本
  ライブラリに渡します。日本語の部分は専用の読み / ふりがなツールで先にかなへ変換し、そのかなを
  渡します（かなはそのまま保持されます）。
- **漢字を事前にかなへ変換する**：かなのチャンクは保持されるため、日本語の漢字を事前にかな
  （ふりがな）へ変換しておけば、誤判定を完全に回避できます。
- その語を多音字として登録しても**解決しません**：多音字の値はピンインであり、やはりピンインの
  経路を通るため、任意の日本語の読みを生成することはできません。

## API リファレンス

| パッケージ | シンボル | 説明 |
| --- | --- | --- |
| `pkg/converter` | `SingleChinesePieceToKana(chineseText string, useHeteronym bool) (string, error)` | テキストをかなに変換します。メインのエントリポイント。 |
| `pkg/converter` | `InitGSEDict(paths ...string) error` | 分かち書きに使う独自の gse 辞書ファイルを 1 つ以上読み込みます。 |
| `pkg/polyphonic` | `LoadPolyphonics(map[string]string)` | 多音字辞書をコンバーターに登録します。 |
| `pkg/polyphonic` | `SafeLoadPolyphonics(map[string]string)` | 上記と同じですが、並行呼び出し向けにミューテックスで保護されています。 |
| `pkg/polyphonic` | `NewManager() *Manager` | メモリ上の多音字マネージャを作成します。既定は `polyphonic.json` と `dict.txt`。 |
| `pkg/polyphonic` | `NewManagerFromFile(targetFilePath, targetDictPath string) (*Manager, error)` | JSON ファイルからマネージャを読み込み、gse 辞書を登録します。 |
| `pkg/polyphonic` | `(*Manager).AddPolyphonic(chinese, pinyin string)` | 多音字を 1 つ追加または上書きします。 |
| `pkg/polyphonic` | `(*Manager).SetTargetFile(targetFilePath, targetDictPath string)` | `Save` / `SaveGSEDict` が使う JSON ファイルと gse 辞書を設定します。 |
| `pkg/polyphonic` | `(*Manager).Save() error` | マネージャを JSON として対象ファイルへ書き出します（既定 `polyphonic.json`）。 |
| `pkg/polyphonic` | `(*Manager).SaveGSEDict() error` | 多音字を gse 辞書ファイルとして書き出します（既定 `dict.txt`）。 |
| `pkg/polyphonic` | `WriteDict(polyphonics map[string]string, targetFilePath string) error` | 中国語 -> ピンインのマップから gse 辞書ファイルを書き出します。 |
| `pkg/polyphonic` | `(*Manager).GetData() map[string]string` | 多音字マップのコピーを返します。 |
| `pkg/polyphonic` | `(*Manager).Initialize()` | 保持しているすべての多音字をコンバーターに登録します。 |

`internal/` 配下のパッケージは実装の詳細であり、公開 API ではありません。

## プロジェクト構成

```
.
├── go.mod                     # module github.com/yukumo-group/Chinese2KanaConverter
├── .goreleaser.yml            # タグ付きリリース用の GoReleaser 設定
├── internal/
│   ├── cpyconverter/          # 中国語 -> ピンイン（go-ego/gpy）、多音字 + gse 辞書
│   ├── japconverter/          # ピンイン -> かな、および非中国語 -> かな
│   ├── language/              # テキストのチャンク分割と中国語判定
│   ├── process/               # 小さな共有ヘルパー
│   └── table/                 # ピンイン / 声母 / 韻母 -> かな のマッピングテーブル
└── pkg/
    ├── converter/             # 公開エントリポイント
    └── polyphonic/            # 多音字 Manager とローダー
```

## 開発

リポジトリをクローンしてテストを実行します：

```bash
git clone https://github.com/yukumo-group/Chinese2KanaConverter.git
cd Chinese2KanaConverter
go mod download
go test ./...
```

CI と同じ方法でテストを実行（詳細出力 + カバレッジ）：

```bash
go test -v ./... -tags=test -coverpkg=./pkg/...,./internal/... -coverprofile=coverage.txt
```

静的解析と lint（CI と同様）：

```bash
go vet ./...
go install golang.org/x/lint/golint@latest
golint -set_exit_status ./...
```

## 継続的インテグレーション

- [`.github/workflows/test.yaml`](./.github/workflows/test.yaml) は push と pull request の
  たびにテスト、`go vet`、`golint` を実行し、カバレッジを Codecov にアップロードします。
- [`.github/workflows/release.yaml`](./.github/workflows/release.yaml) は `v*` にマッチする
  タグが push されたときに GoReleaser を実行します。

## コントリビューション

1. リポジトリを Fork し、機能ブランチを作成します。
2. 既存のスタイルを守ってください：インデントはタブ、エクスポートされたシンボルには
   ドキュメントコメント（golint を通過する必要があります）、テストは `*_test.go` に
   テーブル駆動で書きます。
3. 挙動を変更した場合はテストを追加・更新してください。
4. `go test ./...`、`go vet ./...`、`golint -set_exit_status ./...` がすべて通ることを
   確認します。
5. pull request を作成します。

## ライセンス

[MIT ライセンス](./LICENSE) の下で公開されています。Copyright (c) 2026 Yukumo Group.

