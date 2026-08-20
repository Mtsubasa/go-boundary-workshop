# Ex4: if文から境界値を抽出する

## この章のゴール

`ShippingFee`のbodyをたどり、`if total < 5000`から5000を取り出します。

期待するmessage:

```text
ShippingFee: found boundary value 5000
```

## AST Viewerをもう一度見る

迷ったら次の経路をViewerで確認します。

```text
FuncDecl -> Body -> IfStmt -> Cond(BinaryExpr) -> Y(BasicLit)
```

Coreで対応する条件は`int` parameterが1つの関数と、`parameter < 整数literal`だけです。

## 実装する判定

1. nodeが`*ast.IfStmt`か
2. `Cond`が`*ast.BinaryExpr`か
3. operatorが`token.LSS`か
4. 左辺が関数parameterと同じ名前の`*ast.Ident`か
5. 右辺が整数の`*ast.BasicLit`か

## Hint 1

関数bodyの中は`ast.Inspect(function.Body, func(node ast.Node) bool { ... })`でたどれます。

## Hint 2

型assertionは失敗する可能性があります。失敗したら、そのnodeは今回の対象外として次へ進みます。

```go
ifStatement, ok := node.(*ast.IfStmt)
if !ok {
	return true
}
```

## Hint 3

比較演算子と右辺を確認します。

```go
comparison, ok := ifStatement.Cond.(*ast.BinaryExpr)
if !ok || comparison.Op != token.LSS {
	return true
}

literal, ok := comparison.Y.(*ast.BasicLit)
```

## 完了条件

- `5000`を含む一時診断が表示された
- 診断位置がcondition右辺の5000を指している
- `return 500`の500を境界値として拾っていない

