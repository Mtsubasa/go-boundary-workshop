# Ex5: table testからinputを集める

## この章のゴール

`TestShippingFee`のtableから、`input` fieldの4999と5001だけを集めます。

## 観察するcode

```go
tests := []struct {
	input int
	want  int
}{
	{input: 4999, want: 500},
	{input: 5001, want: 0},
}
```

key付きのfieldは`*ast.KeyValueExpr`です。

```text
KeyValueExpr
├── Key: Ident(input)
└── Value: BasicLit(4999)
```

## 対応づけ

Coreではtest名の`Test`を除いた名前をproduction function名として扱います。

```text
TestShippingFee -> ShippingFee
```

集めた値は次の形で保持します。

```go
map[string][]int64{
	"ShippingFee": {4999, 5001},
}
```

## Hint 1

test file内の`*ast.FuncDecl`から、名前が`Test`で始まる関数だけを対象にします。

## Hint 2

function bodyを`ast.Inspect`し、`*ast.KeyValueExpr`を探します。

## Hint 3

keyを`*ast.Ident`として確認し、`key.Name != "input"`なら収集しません。これがないと`want: 5000`をtest inputと誤認します。

## 完了条件

- `TestShippingFee`に対して4999と5001を取得できた
- `want` fieldの500と0を取得していない

