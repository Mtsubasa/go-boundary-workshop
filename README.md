# go-boundary-workshop

> この`core-complete` branchはStep 3Cの復旧用checkpointです。通常のWorkshopは`main` branchから開始します。この状態を利用する場合は[Recovery Guide](docs/04-RECOVERY.md)の手順で必要なfileだけ復元してください。

GoのASTを使い、境界値のtest case不足を検出する`analysis.Analyzer`を実装するためのrepositoryです。

当日は、[Zenn本](https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop)を**主手順書**として読みながら、このrepositoryのcodeを編集します。

- Zenn本: 背景の説明、実装するcode、完了条件を読む
- このrepository: codeを編集し、commandを実行する
- slide: 全体像と要点を確認する

## このrepositoryで作るもの

`if`文から境界値を抽出し、table-driven testの`input`と照合するCLI toolを作ります。

```go
func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
```

次のtestには境界値の`5000`がないため、完成したanalyzerは警告を出します。

```go
tests := []struct {
	input int
	want  int
}{
	{input: 4999, want: 500},
	{input: 5001, want: 0},
}
```

```text
examples/shipping/shipping.go:4:13: ShippingFee: boundary value 5000 is not tested
```

## repositoryの構成

Workshop中に主に見る場所は次のとおりです。

```text
.
├── boundary/
│   ├── analyzer.go                 # Step 1〜3で実装するfile
│   ├── analyzer_test.go            # 完成したanalyzer用のtest
│   └── testdata/                   # analyzer testの入力codeと期待値一覧
├── cmd/
│   ├── astdump/main.go             # ASTをterminalへ表示するcommand
│   └── boundary/main.go            # 作成するanalyzerのCLI入口
├── examples/shipping/
│   ├── shipping.go                 # analyzerを実行する対象code
│   └── shipping_test.go            # Step 3Cで境界値を追加するtest
├── docs/                            # 仕組みとCIの補足資料
├── go.mod
└── README.md
```

基本的に編集するのは`boundary/analyzer.go`です。Step 3Cだけ、診断が消えることを確認するために`examples/shipping/shipping_test.go`も編集します。

## 始め方

```bash
git clone https://github.com/Mtsubasa/go-boundary-workshop.git
cd go-boundary-workshop
go test ./...
go run ./cmd/boundary ./examples/shipping
```

`main` branchはanalyzerの処理を実装する前のstarterです。そのため、最初の`go run`は診断を表示せずに終了します。これは正常な状態です。

準備ができたら、次の順番で進めます。

## 当日使うcommand

すべてrepositoryのroot directoryで実行します。

| command | 目的 | 実行するタイミング |
|---|---|---|
| `go test ./...` | clone直後の環境確認 | Workshop開始前 |
| `go run ./cmd/astdump "total < 5000"` | Goの式をASTとして表示 | ASTの説明時 |
| `go run ./cmd/boundary ./examples/shipping` | 実装中のanalyzerをexampleへ実行 | 各Stepの完了時 |
| `go test -cover ./examples/shipping` | coverage 100%と境界値不足を比較 | Step 3C |
| `go test -tags=workshop_solution ./boundary` | 完成したanalyzerのtestを実行 | Step 3Cの完了後 |
| `git status --short` | 自分が変更したfileを確認 | 迷ったとき |
| `git diff` | 自分が加えた変更内容を確認 | 迷ったとき |

最も繰り返し使うのは次のcommandです。

```bash
go run ./cmd/boundary ./examples/shipping
```

Stepごとに実装を追加し、このcommandの診断がどう変わるかを確認します。診断を表示した後の`exit status 3`は実装errorではありません。

Step 3Cの完了後に確認するcaseと個別の実行commandは、[`boundary/testdata/README.md`](boundary/testdata/README.md)にまとめています。

## Workshopの進め方

各Stepで、最初にZenn本の対応する章を読み、その後にcodeを編集してcommandを実行します。

| 順番 | Zenn本で読む章 | repositoryで行うこと | 確認command |
|---|---|---|---|
| 0 | ASTを見てみる | `cmd/astdump`で式の構造を見る | `go run ./cmd/astdump "total < 5000"` |
| 1 | Step 1: 関数を見つけて報告する | `boundary/analyzer.go`の`TODO(Step 1)`を実装する | `go run ./cmd/boundary ./examples/shipping` |
| 2 | Step 2: if文から境界値を抽出する | Step 1のcodeを境界値の抽出へ広げる | `go run ./cmd/boundary ./examples/shipping` |
| 3A | Step 3A: table testからinputを集める | production codeとtest codeから値を集める | `go run ./cmd/boundary ./examples/shipping` |
| 3B | Step 3B: 境界値とtest値を分類する | less / boundary / greaterに分類する | `go run ./cmd/boundary ./examples/shipping` |
| 3C | Step 3C: 不足を警告し、修正を確認する | 不足だけを警告し、test caseを追加する | `go run ./cmd/boundary ./examples/shipping` |
| Optional A | analyzer自身をtestする | 用意されたWorkshop用testを実行する | `go test -tags=workshop_solution ./boundary` |
| Optional B | Extra 1: named constantをtypesで評価する | ASTへ型情報を補い、constantを評価する | `go run ./cmd/boundary ./examples/shipping` |
| Optional C | `go vet`・CI | binaryを`go vet`へ組み込み、CIへの接続を確認する | `docs/02-CI_INTEGRATION.md` |

Step 1の入口は[`boundary/analyzer.go`](boundary/analyzer.go)にあります。

```go
// TODO(Step 1): Find function declarations and report their names.
```

Step 2以降は、Zenn本に掲載されたcodeを同じfileへ段階的に追加します。READMEだけを読んで完成codeをコピーするのではなく、各章の説明とASTを確認しながら進めてください。

## 復旧用checkpoint

次のbranchは作業開始地点ではなく、途中で分からなくなった場合に差分を確認するためのcheckpointです。通常は`main` branchのまま作業します。

| branch | 対応する完了地点 |
|---|---|
| `main` | analyzerのstarter |
| `ex03-complete` | Step 1: 関数を検出して報告する実装 |
| `ex04-complete` | Step 2: `if`文から境界値を抽出する実装 |
| `core-complete` | Step 3C: test値を照合して不足を報告する実装 |
| `types-complete` | 発展: named constantとconstant expressionへ対応した実装 |

たとえばStep 1で詰まった場合は、完成checkpointとの差分を確認できます。

```bash
git diff main..ex03-complete -- boundary/analyzer.go
```

作業を壊してしまった場合は、変更をstashへ保存してから、対応するcheckpointの`boundary/analyzer.go`だけを復元できます。

```bash
git stash push -u -m "workshop: before recovery"
git restore --source=ex03-complete -- boundary/analyzer.go
```

Step 1〜Types Extraまでの復元command、元の作業をstashから戻す方法は[Recovery Guide](docs/04-RECOVERY.md)を参照してください。

## ASTを表示する補助command

browserを使わず、式のASTをterminalへ表示できます。

```bash
go run ./cmd/astdump "total < 5000"
```

出力から`BinaryExpr`、`Ident`、`BasicLit`などのnodeを確認できます。

## Documents

1. [Analyzer Architecture](docs/01-ANALYZER_ARCHITECTURE.md): analyzerの処理構造
2. [CI Integration](docs/02-CI_INTEGRATION.md): `go vet`とGitHub Actionsへの組み込み
3. [AST・Types・CI Deep Dive](docs/03-AST_TYPES_CI_DEEP_DIVE.md): ASTからCIまでの仕組みとFAQ
4. [Recovery Guide](docs/04-RECOVERY.md): Stepごとのcheckpointから安全に復元する方法

## 関連リンク

- 実装手順: https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop
- 完成版`go-boundary-checker`: https://github.com/Mtsubasa/go-boundary-checker
