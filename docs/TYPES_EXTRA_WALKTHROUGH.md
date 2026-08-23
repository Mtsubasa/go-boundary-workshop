# Types Extra Walkthrough

この文書は`types-complete` branchの実装解説です。Coreを完了した参加者と、後日Zenn本で学ぶ読者を対象にします。

## 1. ASTだけのCoreが止まる場所

Coreは比較式の右辺を`*ast.BasicLit`として取得します。

```go
if total < 5000 {
```

named constantへ変更すると、右辺は`*ast.Ident`になります。

```go
const freeShippingBoundary = 5000

if total < freeShippingBoundary {
```

ASTから分かるのは、右辺に`freeShippingBoundary`というidentifierが書かれていることまでです。そのidentifierがconstantか、値がいくつかは型検査の結果を使って調べます。

## 2. `TypesInfo.Types`を見る

`analysis.Pass.TypesInfo`には、式と型検査結果の対応が入っています。

```go
typeAndValue, ok := pass.TypesInfo.Types[comparison.Y]
if !ok || typeAndValue.Value == nil {
	return 0, token.NoPos, false
}
```

`TypeAndValue.Value`がnilでなければ、その式はcompile時に評価できるconstantです。これはidentifierだけでなく、次のようなconstant expressionにも使えます。

```go
const freeShippingBoundary = 4000 + 1000
```

## 3. `go/constant`で整数へ変換する

constantの値が`int64`で正確に表現できるかを確認します。

```go
value, exact := constant.Int64Val(typeAndValue.Value)
if !exact {
	return 0, token.NoPos, false
}
```

Coreの数値literalも引き続き扱い、literal以外のconstantだけTypesInfoへ進む構成にしています。

## 4. 一般変数との境界

次の右辺もASTでは同じ`*ast.Ident`です。

```go
var freeShippingBoundary = 5000

if total < freeShippingBoundary {
```

しかし一般変数は実行時に変更される可能性があるため、`TypeAndValue.Value`はnilです。`go/types`はidentifierがどのobjectを指すか、どんな型かを教えますが、実行時の値までは保証しません。

一般変数の値を追うには、代入関係の解析やSSA/data flowが必要です。Extra 1ではここへ踏み込まず、compile時constantまでをsupport contractにします。

## 5. 動作を比べる

named constant不足:

```bash
go run ./cmd/boundary ./examples/shippingconst
```

期待する診断:

```text
ShippingFee: boundary value 5000 is not tested
```

`input: 5000`を追加すると診断が消えます。

自動test:

```bash
go test ./boundary -run TestAnalyzer -v
```

testでは次を区別しています。

- 数値literal: 対応
- named constant: 対応
- constant expression: 対応
- 一般変数: 対象外

## 6. このExtraで伝えること

TypesはASTの代わりではありません。

- ASTで「比較式の右辺」という場所を見つける
- TypesInfoで「その式はcompile時constantか」を確認する
- `go/constant`で「値はいくつか」を取り出す

ASTで構文を見つけ、必要な部分だけ型情報で意味を補う、という役割分担を体験するExerciseです。

