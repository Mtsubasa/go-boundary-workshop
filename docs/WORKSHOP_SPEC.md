# Workshop Core Specification

この文書は、ハンズオン用Core実装の仕様を固定するためのものです。実装中に扱いたい構文が増えても、Coreへ追加せずExtra候補として記録します。

## 1. 学習上の成功条件

参加者が次を体験できれば成功です。

1. AST ViewerでGo sourceとAST nodeを対応づける
2. `analysis.Analyzer`を実行する
3. 自分で書いた処理からsource位置付きの診断を出す
4. production codeとtest codeから集めた値を照合する
5. test caseを追加すると診断が消えることを確認する

AST node名の暗記や、汎用的な静的解析toolの完成は成功条件にしません。

## 2. Golden Path

### Production code

```go
package shipping

func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
```

### 不足があるtest

```go
package shipping

import "testing"

func TestShippingFee(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 4999, want: 500},
		{input: 5001, want: 0},
	}

	for _, tt := range tests {
		if got := ShippingFee(tt.input); got != tt.want {
			t.Errorf("ShippingFee(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
```

### 期待する診断

```text
ShippingFee: boundary value 5000 is not tested
```

### 不足がないtest

次のcaseを追加すると診断が消えます。

```go
{input: 5000, want: 0},
```

## 3. 判定仕様

境界値`b`ごとに、test inputが次の3分類を満たすか確認します。

| 分類 | 条件 | `b = 5000`の代表例 |
|---|---|---:|
| left | `input < b` | 4999 |
| boundary | `input == b` | 5000 |
| right | `input > b` | 5001 |

不足する分類ごとに診断します。

```text
ShippingFee: no test value less than 5000
ShippingFee: boundary value 5000 is not tested
ShippingFee: no test value greater than 5000
```

Golden Pathではleftとrightを用意し、boundaryだけを不足させます。参加者が1 caseを追加するだけで診断が消える体験を優先するためです。

## 4. Coreの入力規約

### Production code

- 対象は`*ast.FuncDecl`
- receiverを持たないpackage-level function
- 引数は1つ
- 引数型は`int`
- 対象条件はfunction body内の`*ast.IfStmt`
- conditionの根は`*ast.BinaryExpr`
- operatorは`token.LSS`
- 左辺はfunction parameterを表すidentifier
- 右辺は整数の`*ast.BasicLit`

### Test code

- production codeと同じpackage
- test function名は`Test<ProductionFunctionName>`
- tableはstruct slice literal
- input fieldのkeyはidentifier`input`
- input fieldのvalueは整数の`*ast.BasicLit`
- `want`など、`input`以外のfieldは収集しない

## 5. Coreで扱わないもの

- `<=`、`>`、`>=`、`==`、`!=`
- `&&`、`||`、括弧を含む複合条件
- 左右を入れ替えた`5000 > total`
- named constantや定数式
- `rate < boundary`のような一般変数
- floatの丸めや隣接値
- 複数parameter
- method
- subtest名からの推測
- positional struct literal
- table field名が`input`以外のtest
- helper functionを通した呼び出し
- external test package
- packageを越えたFacts
- SSAやdata flow

unsupportedな入力を誤って「問題なし」と保証しないよう、READMEへ対応範囲を掲載します。

## 6. Step別の完了条件

### Step 0: ASTを覗く

- `ShippingFee`をAST Viewerへ入力する
- `FuncDecl`、`IfStmt`、`BinaryExpr`、`BasicLit`を画面上で見つける
- node名を暗記させない

### Step 1: 関数を見つけて報告する

- `analysis.Analyzer`をcommandから実行できる
- `FuncDecl`を走査できる
- function名を含む診断がsource位置付きで出る

Step 1用の一時的な診断例:

```text
found function ShippingFee
```

### Step 2: if文から境界値を抽出する

- `FuncDecl`のbodyから`IfStmt`を見つける
- conditionを`BinaryExpr`として確認する
- operatorが`<`であることを確認する
- 右辺の`BasicLit`から`5000`を取得する
- function名、値、位置を保持する

Step 2用の一時的な診断例:

```text
ShippingFee: found boundary value 5000
```

### Step 3: test値と照合する

- `TestShippingFee`から`input` fieldだけを収集する
- production function名とtest function名を対応づける
- left、boundary、rightを判定する
- 不足する分類だけを報告する
- `input: 5000`を追加すると警告が消える

### Analyzerのtest

- 不足があるcaseで期待診断が出る
- 3分類が揃うcaseで診断が出ない
- `want: 5000`をinputとして誤収集しない
- 別のfunctionのtest値を混ぜない

## 7. Extra候補

Core完成後、次の順序で発展させます。

1. named constantを`TypesInfo`と`go/constant`で評価する
2. `&&`と`||`を再帰的にたどる
3. 比較operatorと左右の向きを正規化する
4. identifierのobjectを型情報で識別する
5. test内の実際の`CallExpr`から対象functionを特定する
6. external test packageへFactsで情報を渡す

## 8. Core完成のDefinition of Done

- `go test ./...`が成功する
- Golden Pathでboundary不足の診断が1件出る
- `input: 5000`を追加すると診断が0件になる
- leftだけ不足、rightだけ不足のfixtureがある
- `want` fieldの数値をinputへ混ぜないtestがある
- Coreの対応範囲と非対応範囲がREADMEにある
- 3分以内でstarterから完成動作を説明できる
- 新しいclone先でも同じ結果を再現できる

