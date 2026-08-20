# Ex6: 境界値とtest値を照合する

## この章のゴール

境界値5000とtest inputを照合し、不足している分類だけを診断します。

## 3つに分類する

| 分類 | 条件 | 現在のtest |
|---|---|---|
| less | `input < 5000` | 4999がある |
| boundary | `input == 5000` | ない |
| greater | `input > 5000` | 5001がある |

現在不足しているのはboundaryだけです。

## 実装する

値を順番に見て、3つの`bool`を更新します。

```go
hasLess := false
hasBoundary := false
hasGreater := false
```

走査後、`false`の分類だけ`pass.Reportf`します。

## Hint 1

比較はASTではなく、すでに取り出した`int64`同士に対して行います。

## Hint 2

`switch`の各caseを`input < boundary`、`input == boundary`、`input > boundary`にします。

## Hint 3

boundaryがない場合の診断は次です。

```go
pass.Reportf(
	boundaryPos,
	"%s: boundary value %d is not tested",
	functionName,
	boundaryValue,
)
```

## 動作を確認する

```bash
go run ./cmd/boundary ./examples/shipping
```

期待する診断:

```text
ShippingFee: boundary value 5000 is not tested
```

`examples/shipping/shipping_test.go`へ次を追加します。

```go
{input: 5000, want: 0},
```

もう一度実行し、診断が消えることを確認します。

## 完了条件

- 5000がないときboundary不足を1件だけ報告する
- 5000を追加すると診断が消える

