# Chinese2KanaConverter

[English](./README.md) | **简体中文**

[![Go Test Workflow](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml/badge.svg)](https://github.com/yukumo-group/Chinese2KanaConverter/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/badge/go-1.25.0-00ADD8.svg)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yukumo-group/Chinese2KanaConverter.svg)](https://pkg.go.dev/github.com/yukumo-group/Chinese2KanaConverter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

一个将**中文文本转换为日文片假名（假名）**的 Go 库。

转换器按读音处理中文：先把汉字转成拼音，再把拼音音译为假名。文本中的非中文部分
（英文、已有的假名、标点、数字等）会被单独处理，因此中英日混排的句子会**按块**转换，
而不会被整段丢弃。

```
你好                 ->  ニーハオ
都会区               ->  トウホイチュイ
Hello!               ->  ヘロー！
你好，世界！Hello!    ->  ニーハオシーチエ！ヘロー！
```

> **注意：** 这是面向**日语发音**的*近似音译*，并非对中文的精确转写。声调、送气以及若干
> 音节对立会被有意归并到最接近的假名。详见[发音模型](#发音模型)与[已知限制](#已知限制)。

## 特性

- **中文 -> 假名**：把中文读成拼音，再把每个音节映射为假名。
- **两段式拼音→假名映射**：先查完整的拼音音节表；未命中时退化为把音节拆成*声母* + *韵母*
  分别映射，尽量保证任何拼音输入都能产出结果。
- **区分中英混排**：文本被切分为"中文块 / 非中文块"分别转换。中文块走拼音，其余部分交由
  [English2KanaTransliteration](https://github.com/Luigi-Pizzolito/English2KanaTransliteration) 音译。
- **多音字（异读词）**：可为读音不规则的词（如 `都会区`）登记自定义拼音，既支持直接传入
  `map[string]string`，也支持通过 `Manager` 读写 JSON 文件。
- **并发安全**的字典加载与 `Manager` 访问。

## 工作原理

```
中文文本
     |
     v
SeparateToChunks ......... 切分为「中文 / 非中文」文本块
     |
     +-- 中文 ------> ToPinyin --> ToKana -----------------+
     |                (gpy +       (查表，或                |
     |                 多音字词典)   声母 + 韵母)            |
     |                                                      +--> 假名
     |                                                      |
     +-- 非中文 ----> OthersToKana（英文 / 假名）-----------+
```

1. `language.SeparateToChunks` 逐个字符扫描，把满足 `^[\u4e00-\u9fa5\u3007]+$`（全中文）
   或全非中文的连续片段分组为块。
2. 每个中文块用 `cpyconverter.ToPinyin` 转成拼音。由于声调不影响假名输出，声调会被去除；
   用户登记的多音字词典会先合并进 `gpy` 的短语词典。
3. 每个拼音音节在 `table.PinyinToJapMap` 中查表。若查不到，则改用 `table.InitialToKana`
   与 `table.FinalToKana` 匹配（优先匹配更长的韵母，见 `table.SortStringSlice`）。
4. 每个非中文块交给 `japconverter.OthersToKana`。
5. 各块结果按原顺序拼接。

## 环境要求

- [Go](https://go.dev/dl/) **1.25.0** 或更高版本（见 `go.mod`）。

## 安装

Go module 路径为 `github.com/yukumo-group/Chinese2KanaConverter`（仓库名为
`Chinese2KoeConverter`）。

```bash
go get github.com/yukumo-group/Chinese2KanaConverter
```

## 快速开始

入口函数是 `converter.SingleChinesePieceToKana`。

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
)

func main() {
	// 第二个参数用于开启自定义多音字词典。
	kana, err := converter.SingleChinesePieceToKana("你好", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(kana) // ニーハオ
}
```

输出：

```
ニーハオ
```

第二个参数 `useHeteronym` 决定读拼音时是否查询已登记的多音字词典。

## 中英混排输入

文本按块转换，因此中文与非中文部分可以自由混排：

```go
kana, err := converter.SingleChinesePieceToKana("你好，世界！Hello!", true)
// kana == "ニーハオシーチエ！ヘロー！"
```

每一块单独转换，最后按顺序拼接：

| 文本块 | 类型 | 输出 |
| --- | --- | --- |
| `你好` | 中文 | `ニーハオ` |
| `，` | 非中文 | *（被音译器丢弃）* |
| `世界` | 中文 | `シーチエ` |
| `！` | 非中文 | `！` |
| `Hello` | 非中文 | `ヘロー` |
| `!` | 非中文 | `！` |

## 多音字（异读词）

有些词的读音是"逐字拼音词典"读不对的，可以登记这些读音，让转换器正确取用。

### 直接加载静态映射

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	// 中文词 -> 拼音（音节之间用空格分隔）。
	// 若会在多个 goroutine 中加载，请改用 SafeLoadPolyphonics。
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

### 用 Manager 持久化多音字

`polyphonic.Manager` 持有一个多音字词典，并可以按 JSON 读写。

```go
package main

import (
	"fmt"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

func main() {
	manager := polyphonic.NewManager()

	// 新增一个词，立即对转换器生效。
	manager.AddPolyphonic("都会区", "du hui qu")

	// 写入磁盘……
	manager.SetTargetFile("polyphonic.json")
	if err := manager.Save(); err != nil {
		panic(err)
	}

	// ……之后读回来。
	loaded, err := polyphonic.NewManagerFromFile("polyphonic.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(loaded.GetData()) // map[都会区:du hui qu]

	// 把 Manager 中保存的所有多音字重新登记到转换器。
	loaded.Initialize()
}
```

`Save` 写出的 JSON 形如：

```json
{
  "heteronym": {
    "都会区": "du hui qu"
  }
}
```

> **说明：** `Manager` 内嵌了 `sync.RWMutex`，因此 `NewManager`、`NewManagerFromFile`、
> `GetData` 等方法可安全地并发调用。

## 发音模型

输出是"用日语音系去**近似**中文"，而非精确转写。以下简化是刻意的，直接来自
`internal/table` 中的映射表：

| 方面 | 行为 | 示例 |
| --- | --- | --- |
| 声调 | 直接丢弃（`gpy.Normal`） | 妈 / 麻 / 马 / 骂 -> `マー` |
| 送气 | 不送气与送气塞音共用同一批假名 | `ba` = `pa` = `パー` |
| 卷舌 vs 舌面 | 多个声母合并 | `zhi` / `chi` / `ji` / `qi` -> `チー`；`shi` / `xi` -> `シー` |
| 鼻音韵尾 | `-n` / `-ng` 都变成 `ン` | `yin` = `ying` = `イン` |
| 元音 | 中文 `e` 与 `o` 同样处理 | `e` -> `オー` |
| 音节长度 | 单音节常被加上长音符 | `ta` -> `ター` |

正因如此，不同汉字可能产出同样的假名。若需要更严格的读音，可调整
`internal/table/jap_to_chn.go` 中的映射表，或登记多音字。

## 已知限制

本库的目标是**中文 -> 假名**。英文与日文**假名**借助分块机制可以正常工作，但**日文汉字
不能正确处理**。

| 输入语言 | 是否支持 | 示例 |
| --- | --- | --- |
| 中文 | 支持，但为近似读音（见[发音模型](#发音模型)） | `我爱你` -> `ウォーアイニー` |
| 英文 | 支持 | `Hello` -> `ヘロー` |
| 日文假名（平假名 / 片假名） | 支持 | `こんにちは` -> `コンニチハ` |
| 日文汉字 | **不支持** | `日本語` -> `リーペンユイ`，`東京` -> `トンチン` |

### 为什么日文汉字会读错

`language.IsChinese` 的判定是"每个字符都匹配 `^[\u4e00-\u9fa5\u3007]+$` 时才视为中文"。
而**日文汉字与中文汉字处在同一 Unicode 区块**（U+4E00–U+9FA5），在码位层面无法区分，
于是被路由到拼音路径（`cpyconverter.ToPinyin`），得到的是**中文读音**而非日文读音：

- `日本語` -> `リーペンユイ`（rì běn yǔ），应为 `ニホンゴ`
- `東京` -> `トンチン`（dōng jīng），应为 `トウキョウ`

依赖库 `English2KanaTransliteration` 其实**自带**日文汉字读取能力
（`Kanji_to_kana.go`，数据来自 `dict/kanjidic2_pronounce.json`），并通过
`japconverter.OthersToKana` 接入。但纯汉字的文本块因上述路由**永远到不了它**。另外，
即使走到了那套能力，它的读音是**逐字查表、无上下文**，所以 `日本` 会变成 `ニチホン`
而非 `ニホン`。作为对照，直接调用 `OthersToKana` 得到：

| 输入 | `language.IsChinese` | `OthersToKana` |
| --- | --- | --- |
| `日本語` | `true` | `ニチホンゴ` |
| `東京` | `true` | `トウキョウ` |

### 变通方案

- **调用前先判断语言**：自行按语言切分文本，只把中文部分交给本库；日文部分先用专门的
  读音 / 振假名（furigana）工具转成假名再传入（假名会被原样保留）。
- **预先转换汉字**：由于假名文本块会被保留，先把日文汉字转成假名，即可完全绕开误判。
- 把词登记为多音字**并不能**解决此问题：多音字的值是拼音，仍会走拼音路径，无法产出
  任意的日文读音。

## API 参考

| 包 | 符号 | 说明 |
| --- | --- | --- |
| `pkg/converter` | `SingleChinesePieceToKana(chineseText string, useHeteronym bool) (string, error)` | 将一段文本转换为假名，主入口。 |
| `pkg/polyphonic` | `LoadPolyphonics(map[string]string)` | 向转换器登记一个多音字词典。 |
| `pkg/polyphonic` | `SafeLoadPolyphonics(map[string]string)` | 同上，但带互斥锁，供并发调用者使用。 |
| `pkg/polyphonic` | `NewManager() *Manager` | 创建一个内存版多音字管理器。 |
| `pkg/polyphonic` | `NewManagerFromFile(path string) (*Manager, error)` | 从 JSON 文件加载管理器。 |
| `pkg/polyphonic` | `(*Manager).AddPolyphonic(chinese, pinyin string)` | 新增或覆盖一个多音字。 |
| `pkg/polyphonic` | `(*Manager).SetTargetFile(path string)` | 设置 `Save` 使用的目标文件。 |
| `pkg/polyphonic` | `(*Manager).Save() error` | 以 JSON 形式把管理器写入目标文件。 |
| `pkg/polyphonic` | `(*Manager).GetData() map[string]string` | 返回多音字映射的副本。 |
| `pkg/polyphonic` | `(*Manager).Initialize()` | 把保存的所有多音字登记到转换器。 |

`internal/` 下的包属于实现细节，不属于公开 API。

## 项目结构

```
.
├── go.mod                     # module github.com/yukumo-group/Chinese2KanaConverter
├── .goreleaser.yml            # GoReleaser 发布配置
├── internal/
│   ├── cpyconverter/          # 中文 -> 拼音（go-ego/gpy）+ 多音字词典
│   ├── japconverter/          # 拼音 -> 假名，以及非中文 -> 假名
│   ├── language/              # 文本分块与中文判定
│   ├── process/               # 一些小的公共辅助函数
│   └── table/                 # 拼音 / 声母 / 韵母 -> 假名 映射表
└── pkg/
    ├── converter/             # 公开入口
    └── polyphonic/            # 多音字 Manager 与加载器
```

## 开发

克隆仓库并运行测试：

```bash
git clone https://github.com/yukumo-group/Chinese2KoeConverter.git
cd Chinese2KoeConverter
go mod download
go test ./...
```

与 CI 相同的方式运行测试（详细输出 + 覆盖率）：

```bash
go test -v ./... -tags=test -coverpkg=./pkg/...,./internal/... -coverprofile=coverage.txt
```

静态检查与 lint（与 CI 一致）：

```bash
go vet ./...
go install golang.org/x/lint/golint@latest
golint -set_exit_status ./...
```

## 持续集成

- [`.github/workflows/test.yaml`](./.github/workflows/test.yaml)：在每次 push 与 pull
  request 时运行测试、`go vet` 与 `golint`，并把覆盖率上传到 Codecov。
- [`.github/workflows/release.yaml`](./.github/workflows/release.yaml)：当推送形如 `v*`
  的标签时运行 GoReleaser。

## 贡献

1. Fork 仓库并新建功能分支。
2. 保持现有风格：使用 tab 缩进、为导出符号写文档注释（必须通过 golint），测试写在
   `*_test.go` 中并尽量采用表驱动写法。
3. 任何行为变更都要补充或更新测试。
4. 确保 `go test ./...`、`go vet ./...` 与 `golint -set_exit_status ./...` 全部通过。
5. 提交 pull request。

## 许可证

基于 [MIT 许可证](./LICENSE) 发布。Copyright (c) 2026 Yukumo Group。

