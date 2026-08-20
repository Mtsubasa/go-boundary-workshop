# Ex1: ASTを覗いてみる

## この章のゴール

普段のGo codeとAST nodeが対応していることを目で確かめます。node名を暗記する必要はありません。

## 観察するcode

```go
func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
```

講師が案内するAST Viewerへ上のcodeを貼り付け、source側とtree側を交互に見ます。

## 探すnode

| source | AST node |
|---|---|
| `func ShippingFee` | `FuncDecl` |
| `if` | `IfStmt` |
| `total < 5000` | `BinaryExpr` |
| `5000` | `BasicLit` |

`BinaryExpr`を開き、次の3要素があることを確認します。

- `X`: `total`
- `Op`: `<`
- `Y`: `5000`

## offline fallback

AST Viewerを開けない場合はrepository rootで実行します。

```bash
go run ./cmd/astdump 'total < 5000'
```

## 完了条件

- sourceの`if total < 5000`に対応する`IfStmt`と`BinaryExpr`を見つけた
- `5000`が`BasicLit`であることを確認した

忘れても問題ありません。以降でnodeが必要になったら、またAST Viewerを見ます。

