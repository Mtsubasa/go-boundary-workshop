# 90分Workshop Rehearsal Runbook

この文書は、講義を90分で終えながら「ASTを触るとcodeから自分の診断を出せる」という感動を残すための講師用台本です。全員がCoreを完成することを成功条件にはしません。

## Workshopの成功条件

参加者が次の3つを体験できれば成功です。

1. 短いGo codeとAST Viewerのnodeが対応する
2. 自分で書いた走査処理からsource位置付きの診断が出る
3. 既製toolにない観点も、自作の静的解析で検査できると知る

最低到達点はStep 1の`found function ShippingFee`です。Step 2以降が途中でも、最後の完全版demoで題材の全体像を回収します。

## 前日までに作る講師用環境

live coding用とは別に、完成状態をすぐ表示できるworktreeを用意します。

```bash
git worktree add ../go-boundary-ex03 ex03-complete
git worktree add ../go-boundary-ex04 ex04-complete
git worktree add ../go-boundary-core core-complete
git worktree add ../go-boundary-types types-complete
```

各directoryで次を実行し、dependency downloadを済ませます。

```bash
go test ./...
go run ./cmd/astdump 'total < 5000'
```

さらに、次を準備します。

- Zenn本をbrowserの最初の章で開く
- AST Viewerを`total < 5000`の表示状態で開く
- offline用に`go run ./cmd/astdump 'total < 5000'`を実行済みにする
- slide、editor、terminalの文字sizeを会場後方から読める大きさにする
- `go-boundary-core`のtest tableから5000が抜けた状態を確認する
- CI snippetをslideとZenn本の両方で開けるようにする

## 進行表

| 時間 | 内容 | 参加者の成功表示 | 講師checkpoint |
|---|---|---|---|
| 0〜10分 | coverage 100%への違和感とASTとの出会い | 題材へ納得する | 10分で経験談を閉じる |
| 10〜25分 | ASTの説明、Viewerで式を覗く | `BinaryExpr`を見つける | 20分でViewer操作終了 |
| 25〜30分 | `analysis`の役割 | AST→診断の流れが見える | 30分で必ず実装へ移る |
| 30〜45分 | Step 1: 関数を報告 | `found function ShippingFee` | 42分で一度実行する |
| 45〜60分 | Step 2: 境界値を抽出 | `found boundary value 5000` | 57分で一度実行する |
| 60〜80分 | Step 3: test値と照合 | 5000不足の診断 | 72分で入力収集を切り上げる |
| 80〜90分 | 完全版demo、CI、まとめ | 診断が出て消える | 80分で入力を止める |

## 時間帯別の話し方と操作

### 0〜10分: 出発点を共有する

1. 5000未満で送料が変わるcodeを見せる
2. 4999と5001だけのtestでcoverageが100%になることを見せる
3. 「coverageが間違いなのではなく、測っているものが違う」と説明する
4. 「自分が気にしている観点をcodeにできないか」がtool作成の出発点だったと話す

ここではAIの話を主役にしません。「人でもAIでも、書き手に依存せず同じruleを適用できる」という補助線に留めます。

### 10〜25分: ASTを読む

Viewerへ次の式を入れます。

```go
package sample

func shipping(total int) {
	if total < 5000 {}
}
```

画面で探す順番は`IfStmt`、`BinaryExpr`、`Ident`、`BasicLit`です。node名を暗記させず、「sourceを小さくしてViewerで対応箇所を探す」という調べ方を体験してもらいます。

network不調なら次へ切り替えます。

```bash
go run ./cmd/astdump 'total < 5000'
```

### 25〜30分: analysis frameworkをつなぐ

説明は次の一本道に絞ります。

```text
singlechecker → analysis.Pass → ASTを走査 → pass.Reportf → 診断
```

`analysis.Analyzer`の全fieldやpackage variantの詳細は話しません。`pass.Files`、`pass.Reportf`、`inspect.Analyzer`の3点を、その後に使うものとして紹介します。

### 30〜45分: Step 1

参加者が編集する中心は`boundary/analyzer.go`です。講師は最初に完成形を貼らず、次の順で問いかけます。

1. 関数を表すnodeは何だったか
2. 見つけたいnode typeを`Preorder`へ渡す
3. `FuncDecl.Name.Name`を`pass.Reportf`へ渡す

確認command:

```bash
go run ./cmd/boundary ./examples/shipping
```

期待する診断:

```text
found function ShippingFee
```

42分で未実行なら、残り3分は解説付きで完成codeを提示します。45分時点で全体説明は`ex03-complete`の画面へ切り替えます。

### 45〜60分: Step 2

AST Viewerで見た形を上から狭めます。

```text
IfStmt → Cond → BinaryExpr → Op == token.LSS → Y == BasicLit
```

型assertionが失敗したら`return`する小さなguardを1つずつ追加します。対応範囲を`input < 数値literal`に限定していることを明言します。

期待する診断:

```text
ShippingFee: found boundary value 5000
```

57分で未実行なら、`ex04-complete`の実装を読み、「ASTのどこをcodeに写したか」を確認して次へ進みます。

### 60〜80分: Step 3

最初に完成までの流れを3行で示します。

```text
TestShippingFeeを見つける
input fieldの4999と5001を集める
5000と照合し、不足をReportfする
```

60〜72分はtest値の収集、72〜78分は照合、78〜80分は結果確認に使います。進行が速い参加者にはanalyzer testまたはTypes Extraを案内します。

72分時点で入力収集が未完成なら、以降のlive入力を止め、`core-complete`のcodeを読みながら値の流れを説明します。全員に完成を急がせません。

### 80〜90分: 完全版demo・CI・まとめ

80分になったら進捗に関係なくlive codingを終了します。操作は[Demo Script](DEMO_SCRIPT.md)に従います。

優先順位は次の通りです。

1. coverage 100%を再表示する
2. 自作analyzerが5000不足を報告する
3. 5000を加えると診断が消える
4. `go vet -vettool`またはCI snippetを見せる
5. 3点でまとめる

残り5分未満なら4をslide表示だけにし、まとめの時間を守ります。

## 短縮判断

| 状態 | 判断 | 省くもの |
|---|---|---|
| 30分時点でAST説明中 | 即Step 1へ移る | `analysis`のfield説明 |
| 45分時点でStep 1未完 | `ex03-complete`を読む | live debugの継続 |
| 60分時点でStep 2未完 | `ex04-complete`を読む | literal parseの詳細 |
| 72分時点でtest値未収集 | `core-complete`で説明 | 残りのlive入力 |
| 80分到達 | 完全版demoへ移る | Extraと質疑の一部 |
| 残り5分 | 診断が出て消えるdemoとまとめ | liveの`go vet`実行 |

## 参加者の詰まり方への対応

hintは一度に答えを渡さず、次の順で出します。

1. 調べるnodeを伝える
2. AST Viewerで見るfieldを伝える
3. 該当する型assertionまたは条件式を示す
4. checkpointの完成codeを案内する

声かけ例:

> ここまで動いていれば十分です。続きは一緒に完成版のcodeを読みましょう。先へ進みたい方はそのまま実装を続けて大丈夫です。

## リハーサルの実施方法

本番と同じPC、screen解像度、文字sizeで、説明も省略せず通します。各区間で実時間を記録し、5分以上超えた区間は「話す内容を削る」「入力をstarterへ移す」「checkpointで読む」のいずれかを選びます。

リハーサル後は[Rehearsal Log Template](REHEARSAL_LOG_TEMPLATE.md)を複製し、次回までの修正を3件以内に絞ります。
