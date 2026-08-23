# go-boundary-workshop

Go Conference 2026「こだわりを静的解析で表現しよう 90分で作って動かす自作analyzer入門」のハンズオン用リポジトリです。

> この`main` branchは、参加者が当日作業を始めるstarterです。

## このWorkshopで作るもの

`if`文から境界値を抽出し、境界値そのものと境界の左右に属する値がtable-driven testに含まれていなければ警告するCLI toolを作ります。

```go
func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
```

次のtestには`5000`がないため、analyzerが警告します。

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
ShippingFee: boundary value 5000 is not tested
```

## 事前確認

codeを変更する前に、すべてのtestが成功することを確認します。

```bash
go test ./...
```

次に、まだ診断を実装していないstarterのanalyzerを起動します。

```bash
go run ./cmd/boundary ./examples/shipping
```

何も表示されず終了すれば準備完了です。

## Workshopで確認する違和感

題材の通常testとcoverageを確認します。

```bash
go test -cover ./examples/shipping
```

境界値5000は不足していますが、両方の分岐を通るためstatement coverageは100%です。

```text
coverage: 100.0% of statements
```

WorkshopのStep 3まで完成したanalyzerを実行すると、coverageには現れなかった不足を報告します。

```bash
go run ./cmd/boundary ./examples/shipping
```

`examples/shipping/shipping_test.go`には`input: 5000`がないため、完成後は次の診断と終了status 3が表示されます。

```text
examples/shipping/shipping.go:4:13: ShippingFee: boundary value 5000 is not tested
exit status 3
```

test tableへ`{input: 5000, want: 0}`を追加してもう一度実行すると、診断が消えます。

この「coverageは100%、それでも検査したい観点が残る」という差が、本Workshopの出発点です。当日の完全版demo手順は[Demo Script](docs/DEMO_SCRIPT.md)にまとめています。

## 3 Step

1. 関数を見つけて報告する
2. `if`文から境界値を抽出する
3. test値と照合して、不足しているcaseを警告する

## ASTを覗くoffline fallback

当日はbrowser上のAST Viewerを使います。接続できない場合は、標準libraryだけで動く補助commandでも同じnodeを確認できます。

```bash
go run ./cmd/astdump 'total < 5000'
```

出力から`BinaryExpr`、`Ident`、operatorの`<`、`BasicLit`を探します。node名を暗記する必要はありません。実装ではoperatorの`<`を`token.LSS`として判定します。

## Coreの対応範囲

90分で静的解析が初めての参加者も完成できるよう、Coreでは次のcodeだけを対象とします。

- 引数が1つの関数
- 引数と境界値は`int`
- 条件は`if input < 数値リテラル`
- 対応演算子は`<`のみ
- test名は`Test<関数名>`
- tableの入力field名は`input`
- test値は数値リテラル
- production codeとtest codeは同じpackage

複合条件、named constant、変数境界、他の比較演算子、複数引数、external test packageはExtraまたは完成版で扱います。

詳細は次の文書を参照してください。

- [Workshop Spec](docs/WORKSHOP_SPEC.md): Coreの対応範囲と成功条件
- [Exercise Drafts](docs/exercises/README.md): Zenn本へ展開する参加者向け手順
- [Core Walkthrough](docs/CORE_WALKTHROUGH.md): 実装と切り出しの講師用解説
- [Demo Script](docs/DEMO_SCRIPT.md): 80〜90分の完全版demo台本

## 参加前の準備

- Go codeを編集・実行でき、Gitを利用できるノートPC
- ノートPCの電源adapter

開催前に、次が成功することを確認してください。

```bash
git clone <repository-url>
cd go-boundary-workshop
go test ./...
```

Gitの詳しい操作、GitHub account、ASTや静的解析の事前知識は必要ありません。

現在の完成版はGo 1.25.12で動作確認しています。必要な最低versionは`go.mod`を参照してください。

## 教材

- Zenn本: 公開後にlinkを追加
- 登壇資料: 公開後にlinkを追加
- 完成版`go-boundary-checker`: https://github.com/Mtsubasa/go-boundary-checker

## checkpoint

講師による復旧と、実装を先に確認したい場合のために次のbranchを用意します。

- `main`: 参加者が開始するstarter
- `ex03-complete`: 関数を見つけて報告できるcheckpoint
- `ex04-complete`: 境界値を抽出できるcheckpoint
- `core-complete`: test値との照合まで完成したcheckpoint

参加者にはbranch切り替えを要求せず、checkpointは講師による復旧と進行の速い参加者向けに使います。
