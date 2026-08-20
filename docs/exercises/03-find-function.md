# Ex3: 関数を見つけて報告する

## この章のゴール

`*ast.FuncDecl`を見つけ、関数名をsource位置付きの診断として表示します。

期待するmessage:

```text
found function ShippingFee
```

## 編集する場所

`boundary/analyzer.go`の`run`にあるEx3のTODOを編集します。starterには`inspect.Analyzer`と`inspector.Inspector`を取得する処理まで用意します。

## 考えること

1. どの種類のnodeだけを受け取りたいか
2. 受け取った`ast.Node`を`*ast.FuncDecl`として扱えるか
3. 関数名とsource位置をどのfieldから取るか

## 実行する

```bash
go run ./cmd/boundary ./examples/shipping
```

## Hint 1

`inspector.Inspector.Preorder`へ、対象nodeのfilterとcallbackを渡します。

## Hint 2

node filterは次の形です。

```go
[]ast.Node{(*ast.FuncDecl)(nil)}
```

## Hint 3

callbackで受け取ったnodeを型assertionし、`Name.Name`と`Name.Pos()`を使います。

```go
function := node.(*ast.FuncDecl)
pass.Reportf(function.Name.Pos(), "found function %s", function.Name.Name)
```

## 完了条件

- `ShippingFee`を含む診断が表示された
- errorのfile名と行番号が`shipping.go`を指している

この時点で、自分で書いた処理がGo sourceに対する診断になりました。

