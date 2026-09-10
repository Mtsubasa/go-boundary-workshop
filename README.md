# go-boundary-workshop

> このブランチは応用課題3の完了状態です。通常は`main`ブランチから開始し、名前付き定数への対応で詰まった場合の確認・復旧用として利用してください。

GoのASTを使い、境界値のテストケース不足を検出する`analysis.Analyzer`を実装するためのリポジトリです。

当日は、[Zenn本](https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop)を**主な手順書**として読みながら、このリポジトリのコードを編集します。

- Zenn本：背景と実装手順を読む
- このリポジトリ：コードを編集し、コマンドを実行する
- スライド：全体像と要点を確認する

## このリポジトリで作るもの

`if`文から境界値を抽出し、テーブルテストの`input`と照合するCLIツールを作ります。

```go
func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
```

次のテストには境界値の`5000`がないため、完成したanalyzerは診断を出します。

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

## リポジトリの構成

ワークショップ中に主に見る場所は次のとおりです。

```text
.
├── boundary/
│   ├── analyzer.go                 # Step 1〜3で実装するファイル
│   ├── analyzer_workshop_test.go   # 完成したanalyzer用のテスト
│   └── testdata/                   # analyzerテストの入力コードと期待値一覧
├── cmd/
│   ├── astdump/main.go             # ASTをターミナルへ表示するコマンド
│   └── boundary/main.go            # 作成するanalyzerのCLI入口
├── examples/shipping/
│   ├── shipping.go                 # analyzerを実行する対象コード
│   └── shipping_test.go            # Step 3Cで境界値を追加するテスト
├── docs/                            # 構成・CI・復旧方法の補足資料
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

`main`ブランチはanalyzerの処理を実装する前の初期状態です。そのため、最初の`go run`は何も表示せずに終了します。これは正常な状態です。

準備ができたら、次の順番で進めます。

## 当日使うコマンド

すべてリポジトリのルートで実行します。

| コマンド | 目的 | 実行するタイミング |
|---|---|---|
| `go test ./...` | clone直後の環境確認 | ワークショップ開始前 |
| `go run ./cmd/astdump "total < 5000"` | Goの式をASTとして表示 | ASTの説明時 |
| `go run ./cmd/boundary ./examples/shipping` | 実装中のanalyzerをサンプルコードへ実行 | 各Stepの完了時 |
| `go test -cover ./examples/shipping` | カバレッジ100%と境界値不足を比較 | Step 3C |
| `go test -tags=workshop_solution ./boundary` | 完成したanalyzerのテストを実行 | Step 3Cの完了後 |
| `git status --short` | 自分が変更したファイルを確認 | 迷ったとき |
| `git diff` | 自分が加えた変更内容を確認 | 迷ったとき |

最も繰り返し使うのは次のコマンドです。

```bash
go run ./cmd/boundary ./examples/shipping
```

Stepごとに実装を追加し、このコマンドの表示がどう変わるかを確認します。`pass.Reportf`でメッセージや診断を表示した後の`exit status 3`は、実装エラーではありません。

Step 3Cの完了後に確認するケースと個別の実行コマンドは、[`boundary/testdata/README.md`](boundary/testdata/README.md)にまとめています。

## ワークショップの進め方

各Stepでは、最初にZenn本の対応する章を読み、その後にコードを編集してコマンドを実行します。

| 順番 | Zenn本で読む章 | リポジトリで行うこと | 確認コマンド |
|---|---|---|---|
| 0 | ASTを見てみる | `cmd/astdump`で式の構造を見る | `go run ./cmd/astdump "total < 5000"` |
| 1 | Step 1：関数を見つけて報告する | `boundary/analyzer.go`の`TODO(Step 1)`を実装する | `go run ./cmd/boundary ./examples/shipping` |
| 2 | Step 2：if文から境界値を抽出する | Step 1のコードを境界値の抽出へ広げる | `go run ./cmd/boundary ./examples/shipping` |
| 3A | Step 3A：テーブルテストから入力値を集める | 実装コードとテストコードから値を集める | `go run ./cmd/boundary ./examples/shipping` |
| 3B | Step 3B：境界値とテスト値を分類する | 未満・境界値・超過に分類する | `go run ./cmd/boundary ./examples/shipping` |
| 3C | Step 3C：不足を診断し、修正を確認する | 不足だけを診断し、テストケースを追加する | `go run ./cmd/boundary ./examples/shipping` |
| 応用課題1 | analyzer自身をテストする | 用意されたテストを実行する | `go test -tags=workshop_solution ./boundary` |
| 応用課題2 | 比較式の左右を正規化する | `total < 5000`と`5000 > total`へ対応する | `go run ./cmd/boundary ./examples/shipping` |
| 応用課題3 | 名前付き定数を型情報で評価する | ASTへ型情報を補い、定数を評価する | `go run ./cmd/boundary ./examples/shipping` |
| 応用課題4 | 成功メッセージを表示する | 診断の有無を記録し、問題がなければ結果を表示する | `go run ./cmd/boundary ./examples/shipping` |
| 参考 | `go vet`・CI | 実行ファイルを`go vet`へ組み込み、CIへの接続を確認する | `docs/02-CI_INTEGRATION.md` |

4つの応用課題は互いに独立しています。残り時間や興味に合わせて、好きなものを1つ選びます。応用課題1は用意されたテストを実行するだけ、応用課題2は比較式の左右へ対応する実装、応用課題3は型情報を扱う難しめの実装、応用課題4は実行結果を分かりやすくする小さな実装です。

Step 1の入口は[`boundary/analyzer.go`](boundary/analyzer.go)にあります。

```go
// TODO(Step 1): 関数宣言を探し、関数名を表示する。
```

Step 2以降は、Zenn本に掲載されたコードを同じファイルへ段階的に追加します。READMEだけを読んで完成コードをコピーするのではなく、各章の説明とASTを確認しながら進めてください。

## 復旧用チェックポイント

次のブランチは作業開始地点ではなく、途中で分からなくなった場合に差分を確認するためのチェックポイントです。通常は`main`ブランチのまま作業します。

| ブランチ | 対応する完了地点 |
|---|---|
| `main` | analyzerの初期状態 |
| `ex03-complete` | Step 1: 関数を検出して報告する実装 |
| `ex04-complete` | Step 2: `if`文から境界値を抽出する実装 |
| `step3-complete` | Step 3C：テスト値を照合して不足を診断する実装 |
| `types-complete` | 応用課題3：名前付き定数と定数式へ対応した実装 |

たとえばStep 1で詰まった場合は、完成チェックポイントとの差分を確認できます。

```bash
git diff main..ex03-complete -- boundary/analyzer.go
```

作業を壊してしまった場合は、変更をstashへ保存してから、対応するチェックポイントの`boundary/analyzer.go`だけを復元できます。

```bash
git stash push -u -m "workshop: before recovery"
git restore --source=ex03-complete -- boundary/analyzer.go
```

Step 1から応用課題3までの復元コマンドと、元の作業をstashから戻す方法は[復旧手順](docs/03-RECOVERY.md)を参照してください。

## ASTを表示する補助コマンド

ブラウザを使わず、式のASTをターミナルへ表示できます。

```bash
go run ./cmd/astdump "total < 5000"
```

出力から`BinaryExpr`、`Ident`、`BasicLit`などのノードを確認できます。

## 補足資料

1. [Analyzerの構成](docs/01-ANALYZER_ARCHITECTURE.md)：analyzerの処理構造
2. [CIへの組み込み](docs/02-CI_INTEGRATION.md)：`go vet`とGitHub Actionsへの組み込み
3. [復旧手順](docs/03-RECOVERY.md)：Stepごとのチェックポイントから安全に復元する方法

## 関連リンク

- 実装手順: https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop
- 完成版`go-boundary-checker`: https://github.com/Mtsubasa/go-boundary-checker
