# go-boundary-workshop

GoのASTを使い、境界値のtest case不足を検出する`analysis.Analyzer`を実装するためのrepositoryです。

詳しい実装手順は[Zenn本](https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop)にまとめています。このrepositoryでは、実行するcodeと技術資料を管理します。

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

## 始め方

```bash
git clone https://github.com/Mtsubasa/go-boundary-workshop.git
cd go-boundary-workshop
go test ./...
go run ./cmd/boundary ./examples/shipping
```

`main` branchはanalyzerの処理を実装する前のstarterです。そのため、最初の`go run`は診断を表示せずに終了します。実装はZenn本の手順に沿って進められます。

## 実装状態

| branch | 内容 |
|---|---|
| `main` | analyzerのstarter |
| `ex03-complete` | 関数を検出して報告する実装 |
| `ex04-complete` | `if`文から境界値を抽出する実装 |
| `core-complete` | test値を照合して不足を報告する実装 |
| `types-complete` | named constantとconstant expressionへ対応した実装 |

## ASTを表示する補助command

browserを使わず、式のASTをterminalへ表示できます。

```bash
go run ./cmd/astdump "total < 5000"
```

出力から`BinaryExpr`、`Ident`、`BasicLit`などのnodeを確認できます。

## Documents

1. [Analyzer Architecture](docs/01-ANALYZER_ARCHITECTURE.md): analyzerの処理構造
2. [CI Integration](docs/02-CI_INTEGRATION.md): `go vet`とGitHub Actionsへの組み込み

## 関連リンク

- 実装手順: https://zenn.dev/tsubasa_m/books/go-static-analysis-workshop
- 完成版`go-boundary-checker`: https://github.com/Mtsubasa/go-boundary-checker
