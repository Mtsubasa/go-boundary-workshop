# Extra 1: named constantをtypesで評価する

このExerciseはCoreを早く終えた人と、後日教材を読む人向けです。

## Coreを壊す

境界値をliteralからnamed constantへ変更します。

```go
const freeShippingBoundary = 5000

func ShippingFee(total int) int {
	if total < freeShippingBoundary {
		return 500
	}
	return 0
}
```

Coreは右辺が`*ast.BasicLit`であることを前提としているため、境界値を取得できません。AST上の右辺は`*ast.Ident`です。

## この章のゴール

`analysis.Pass.TypesInfo`と`go/constant`を使い、右辺がliteralでもnamed constantでも整数値5000を得られるようにします。

## ASTとtypesの役割

| 情報 | AST | types |
|---|---:|---:|
| 右辺がidentifierである | 分かる | 分かる |
| identifierの名前 | 分かる | 分かる |
| constantを指している | 分からない | 分かる |
| constantの値が5000 | 分からない | 分かる |

## Hint 1

`pass.TypesInfo.Types[comparison.Y]`から、その式の`TypeAndValue`を取得します。

## Hint 2

compile時定数であれば`TypeAndValue.Value`はnilではありません。

## Hint 3

整数として表現できるかは`go/constant.Int64Val`で確認できます。

```go
typeAndValue, ok := pass.TypesInfo.Types[comparison.Y]
if !ok || typeAndValue.Value == nil {
	// compile時に値が決まらない式
}

value, exact := constant.Int64Val(typeAndValue.Value)
```

## 完了条件

- `total < 5000`を引き続き扱える
- `total < freeShippingBoundary`から5000を取得できる
- 実行時に決まる一般変数を「typesだけで値が分かる」と扱わない

一般変数の値を知るには代入の追跡やSSA/data flowが必要です。ここでは実装対象をcompile時定数までに限定します。
