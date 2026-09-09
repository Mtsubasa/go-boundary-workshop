# 03. AST・Types・CI Deep Dive

この文書は、Workshop参加者が各APIのつながりと実装の背景を振り返るための補足資料です。Workshop中にすべて読む必要はありません。手順を進めた後で、「なぜこのAPIを使うのか」「どこまで保証できるのか」を深掘りしたいときに参照してください。

## 全体像

今回のtoolは、次の層を順につないでいます。

```text
Go source code
  ↓ parse
AST（codeの構文）
  ↓ analysis.Pass
Analyzer（探すruleと診断）
  ↓ 必要な場合だけ参照
TypesInfo（identifierや式の意味）
  ↓ driverで実行
CLI / go vet -vettool
  ↓ 終了status
CI
```

重要なのは、それぞれの役割を混ぜないことです。

| 層 | 答えられること | 答えられないこと |
|---|---|---|
| AST | どんな構文が、どこに書かれているか | identifierが実際に何を指すか、実行時の値 |
| TypesInfo | 式の型、identifierの参照先、compile時constantの値 | 一般変数へ実行時に入る値 |
| Analyzer | 何を検査して、どこへ診断を出すか | どの環境・タイミングで実行するか |
| driver / CI | Analyzerをpackageへ適用し、結果を成功・失敗へ変える | ruleそのものの意味 |

## 1. ASTはsource codeの何を表すのか

ASTはAbstract Syntax Treeの略で、source codeを文法上のまとまりへ分解したtreeです。空白や字下げではなく、宣言、文、式といった構造を表します。

```go
if total < 5000 {
	return 500
}
```

今回見る範囲では、次のようなtreeになります。

```text
IfStmt
├── Cond: BinaryExpr
│   ├── X: Ident(total)
│   ├── Op: token.LSS
│   └── Y: BasicLit(5000)
└── Body: BlockStmt
    └── ReturnStmt
        └── BasicLit(500)
```

同じ`500`という文字列でも、条件の右辺にあるか、`return`の値にあるかでtree上の位置が違います。そのため、文字列検索ではなく「`IfStmt.Cond`の`BinaryExpr.Y`」をたどることで、境界値だけを選べます。

### `ast.Node`とは

ASTを構成する型は`ast.Node` interfaceを実装し、少なくとも`Pos()`と`End()`を持ちます。nodeは大きく次の種類に分かれます。

| 種類 | 例 | 今回の例 |
|---|---|---|
| declaration | `FuncDecl`、`GenDecl` | `func ShippingFee...` |
| statement | `IfStmt`、`ReturnStmt` | `if total < 5000` |
| expression | `BinaryExpr`、`Ident`、`BasicLit` | `total < 5000`、`total`、`5000` |

`*ast.FuncDecl`のようにpointer型で扱うのは、tree上のnode自体を参照し、位置情報やfieldを読むためです。

### `token.Pos`は行番号ではない

`token.Pos`は、それ単体ではfile名や行番号ではなく、`token.FileSet`の中で管理される位置です。`analysis.Pass`では次の2つが対応しています。

```go
pass.Fset                         // source位置の管理表
pass.Files                        // 同じ位置表を使うAST
pass.Fset.Position(node.Pos())    // file名・行・列へ変換
```

`pass.Reportf(node.Pos(), ...)`へ位置を渡すと、driverが`FileSet`を使い、人が読める`file:line:column`へ変換します。

### なぜ正規表現ではないのか

正規表現でも`< 5000`という文字列は探せます。しかし、comment、文字列literal、`return 5000`、別の変数との比較などを区別する処理が急速に複雑になります。ASTならGo parserが文法を解釈済みなので、ruleは必要な構造だけに集中できます。

### `ast.Inspect`と`inspector.Inspector`の違い

どちらもASTを走査しますが、用途が少し違います。

| API | 向いている用途 | 今回の使い方 |
|---|---|---|
| `ast.Inspect` | あるnode以下を手軽に深さ優先でたどる | 1つの関数bodyから`IfStmt`を探す |
| `inspector.Inspector` | 複数fileを、node型で絞って走査する | package全体から`FuncDecl`を探す |

今回のcodeでは、外側を`Inspector.Preorder`で関数に絞り、内側を`ast.Inspect`で関数bodyだけたどります。

## 2. `go/analysis`は何をしているのか

`go/analysis`では、検査ruleと実行方法を分けます。

```text
Analyzer = ruleの定義
Pass     = 1つのpackageを検査するときの入力と出力
driver   = packageを読み込み、Analyzerを実行するprogram
```

### `analysis.Analyzer`

```go
var Analyzer = &analysis.Analyzer{
	Name: "boundary",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}
```

`Analyzer`自体はASTを持っていません。「名前」「説明」「実行する関数」「先に必要な解析」を宣言する、ほぼ不変の設定値です。

### `analysis.Pass`

`Run`はpackageごとに`Pass`を受け取ります。今回よく使うfieldは次のとおりです。

| field / method | 役割 |
|---|---|
| `pass.Files` | packageに含まれるGo fileのAST |
| `pass.Fset` | AST nodeをsource位置へ対応づける |
| `pass.Pkg` | 型検査済みpackage情報 |
| `pass.TypesInfo` | 式やidentifierに対応する型情報 |
| `pass.ResultOf` | `Requires`で指定したAnalyzerの結果 |
| `pass.Reportf` | source位置付きの診断を報告 |

`Run`が返す`error`と`pass.Reportf`は意味が違います。

- `Reportf`: 検査対象のcodeにrule違反を見つけた
- `error`: Analyzer自体が処理を続けられなかった

境界値caseがないことはtoolの故障ではないため、`error`ではなく診断として報告します。

### `Requires`と`ResultOf`

```go
Requires: []*analysis.Analyzer{inspect.Analyzer}
```

この宣言により、driverは先に`inspect.Analyzer`を実行します。その結果は次のように受け取れます。

```go
inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
```

自分で各fileから`Inspector`を作らないのは、依存する解析の準備と再利用をframeworkへ任せるためです。

### `singlechecker`は何をしているのか

```go
func main() {
	singlechecker.Main(boundary.Analyzer)
}
```

`singlechecker`は1つのAnalyzerを実行するdriverです。通常のCLIとしてpackage patternを受け取れるだけでなく、`go vet -vettool`から渡される設定fileも認識し、内部で`unitchecker`として動作します。

つまり、`boundary.Analyzer`のruleを書き換えずに、次の両方で使えます。

```bash
go run ./cmd/boundary ./examples/shipping
go vet -vettool="$(pwd)/boundary-checker" ./examples/shipping
```

### なぜ同じ関数を二重に見ることがあるのか

testを含むpackageを解析すると、driverは通常fileだけのpackageと、test fileを含むvariantを扱うことがあります。今回のruleはproduction codeとtest codeを同時に比較するため、test fileを含むvariantだけを利用します。

これはGoのすべてのpackage loadingを学ぶための処理ではなく、Workshopのruleに必要な入力を1つの`Pass`へ揃えるための処理です。

## 3. ASTへTypesInfoを足す理由

次の2つは、見た目の目的は同じでもAST nodeが違います。

```go
if total < 5000                 // Yは*ast.BasicLit
if total < freeShippingBoundary // Yは*ast.Ident
```

ASTから`freeShippingBoundary`という名前は取得できます。しかし、その名前がconstantなのか、別scopeの変数なのか、値がいくつなのかは名前だけでは判断できません。

### 型情報はdriverが準備済み

`analysis.Pass.TypesInfo`には、driverがpackageを型検査した結果が入っています。Analyzer側で`types.Config.Check`をもう一度呼ぶ必要はありません。

今回使うのは`Types` mapです。

```go
typeAndValue, ok := pass.TypesInfo.Types[comparison.Y]
```

`TypeAndValue`の主なfieldは次の2つです。

| field | 内容 |
|---|---|
| `Type` | 式の型 |
| `Value` | compile時constantなら、その値。それ以外はnil |

### literalもconstantも「式」

Goでは次の式はcompile時に値を確定できます。

```go
5000
freeShippingBoundary
4000 + 1000
```

`TypesInfo.Types[expression].Value`を使うと、構文の形が違ってもconstant valueとして共通に扱えます。`go/constant.Int64Val`は、その値を`int64`で正確に表現できるかも同時に確認します。

```go
value, exact := constant.Int64Val(typeAndValue.Value)
if !exact {
	return 0, token.NoPos, false
}
```

`exact`を無視しないのは、範囲外の整数を誤った値へ丸めて診断しないためです。

### `Types`、`Uses`、`Defs`の使い分け

混同しやすい3つは、それぞれ答える質問が違います。

| map | 質問 |
|---|---|
| `Types[expr]` | この式の型とconstant valueは何か |
| `Uses[id]` | この使用箇所のidentifierは、どのobjectを参照するか |
| `Defs[id]` | この宣言箇所のidentifierは、どのobjectを定義するか |

今回は「右辺の式を整数constantとして評価できるか」が知りたいため、`Types`が最短です。同名変数の区別や、parameterと比較左辺が同じobjectかを厳密に確認するなら`Uses`と`Defs`を使います。

### 一般変数の値はなぜ分からないのか

```go
var boundary = 5000
boundary = loadFromConfig()
if total < boundary {
```

型検査で`boundary`が`int`だとは分かりますが、実行時の値は制御経路や代入によって変わります。その追跡にはSSA、control-flow、data-flowなど別の解析が必要です。

Types Extraがcompile時constantまでを対象にするのは、難しい処理を省略しただけではなく、toolが確実に保証できる範囲を定めた設計です。

## 4. AnalyzerをCIへつなぐ

CIへ新しい解析logicを書くわけではありません。localで動いたAnalyzerをbinaryにし、同じcommandをCIで実行します。

```text
boundary.Analyzer
  ↓ singlechecker.Main
boundary-checker binary
  ↓ go vet -vettool
終了status 0 / 非0
  ↓
GitHub Actionsの成功 / 失敗
```

### localで先に確認する理由

```bash
go build -o boundary-checker ./cmd/boundary
go vet -vettool="$(pwd)/boundary-checker" ./examples/shipping
```

localで失敗するcommandを先に直しておけば、CIで失敗したときにworkflowの問題とAnalyzerの問題を分けて考えられます。

`-vettool`へ渡せるのは、`go vet`とのprotocolに対応したanalysis driverです。今回の`singlechecker`製binaryは、`go vet`から起動された場合に`unitchecker`として動作するため利用できます。

### `go test`とは別stepにする理由

```yaml
- name: Run tests
  run: go test ./...

- name: Check boundary test cases
  run: go vet -vettool="${RUNNER_TEMP}/boundary-checker" ./...
```

両者が検査している内容は異なります。

| command | 主な問い |
|---|---|
| `go test` | 実行したcaseで期待する結果になったか |
| boundary analyzer | teamが必要とする境界値caseがsource上に存在するか |

stepを分けると、CI画面だけで失敗理由を判別できます。

### `./...`と`testdata`

Go commandのpackage patternでは、`testdata`という名前のdirectoryは通常のpackage探索から除外されます。そのため、CIの`./...`が`boundary/testdata/*`を直接解析するわけではありません。

`boundary/testdata`は`boundary` packageのanalyzer testが明示的にfileを読み、Analyzerへ渡すfixtureです。実projectを検査する`go vet`と、Analyzer自身を検査するfixture testは別の経路です。

### CIで診断が出たら何が起こるか

`go vet`は診断があると非0で終了し、GitHub Actionsのstepが失敗します。これは「Analyzerがcrashした」という意味ではなく、設定したruleに違反するcodeが見つかったという意味です。

運用では、次を決めておく必要があります。

1. どのpackage patternを対象にするか
2. どの構文をsupport対象とするか
3. 診断を必須checkとして扱うか
4. false positiveが出た場合の修正・除外方法

## 5. よくある質問への短い回答

### ASTはcompilerそのものですか

ASTはcompilerも利用する構文表現ですが、今回作るAnalyzerはcompilerではありません。driverがparse・型検査した結果を受け取り、独自ruleで診断します。

### ASTを変更するとsource codeも書き換わりますか

今回のAnalyzerはASTを読むだけなので書き換えません。修正案を扱う`SuggestedFix`もありますが、Coreでは使用しません。

### node名を全部覚える必要がありますか

ありません。最小のsource codeをAST Viewerへ入れ、必要なfieldだけpackage documentで確認するのが実用的です。

### `Ident.Name`を見るだけでは駄目ですか

構文を限定したCoreでは十分な場合があります。しかし、同じ名前が別scopeに存在する可能性まで区別するなら`TypesInfo.Uses`や`Defs`でobjectを比較します。

### `Reportf`と`fmt.Printf`は何が違いますか

`Reportf`は位置付きの構造化された診断としてdriverへ渡ります。driverはCLI表示、JSON、`go vet`の終了statusなど実行環境に応じた形へ変換できます。

### なぜ診断時に`exit status 3`が出るのですか

standalone driverが診断ありを非0の終了statusで表すためです。Workshopの実行例では3になります。shellやCIからは「検査に通らなかった」と判定できます。

### CIは変更したfileだけを検査しますか

このworkflow例では違います。指定したpackage pattern全体を毎回検査します。差分だけに限定するには、変更fileからpackageを選ぶ別の仕組みが必要です。

### Typesを使えば実行時の値も分かりますか

分かりません。型、参照先、compile時constantなどは分かりますが、一般変数の実行時値にはdata-flow解析などが必要です。

### このAnalyzerは境界値testを完全に保証しますか

保証するのはsupport matrixに書いた構文だけです。`<=`、複合条件、methodなどを対象外にしているため、あらゆるGo codeの境界値testを保証するtoolではありません。

## 6. 理解確認

次の質問へ自分の言葉で答えると、このWorkshopで扱った仕組みを整理できます。

1. `return 500`ではなく`if`右辺の5000だけを選べるのはなぜか
2. `token.Pos`からfile名と行番号が出るまでに何が必要か
3. `Analyzer`、`Pass`、`singlechecker`の役割はどう違うか
4. `Requires`へ`inspect.Analyzer`を書くと何が起こるか
5. 診断を`error`として返さず`Reportf`するのはなぜか
6. named constantがASTだけでは評価できないのはなぜか
7. constantと一般変数を`TypesInfo`でどう見分けるか
8. `go test`とcustom analyzerをCIで別stepにするのはなぜか
9. `go vet ./...`が`testdata`を直接検査しないのはなぜか
10. 今回のAnalyzerが保証しない構文を3つ説明できるか

## 参考資料

- [`go/ast`](https://pkg.go.dev/go/ast)
- [`go/token`](https://pkg.go.dev/go/token)
- [`go/types`](https://pkg.go.dev/go/types)
- [`go/constant`](https://pkg.go.dev/go/constant)
- [`analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [`singlechecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/singlechecker)
- [`unitchecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/unitchecker)
- [`go vet`](https://pkg.go.dev/cmd/go#hdr-Report_likely_mistakes_in_packages)
