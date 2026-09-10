# Analyzerのテストデータ

このディレクトリには、完成したanalyzer自身をテストするための入力コードを置いています。

`examples/shipping`はワークショップ中に編集して動きを確認する題材です。一方、ここにある各ディレクトリは編集せず、analyzerが期待どおりの診断を出すか確認するために使います。

## 今回確認するケース

| ディレクトリ | 確認すること | 期待する診断 |
|---|---|---|
| `missing_boundary` | 境界の左右はあるが、境界値そのものがない | `ShippingFee: boundary value 5000 is not tested` |
| `complete` | 境界値の未満・境界値・超過がすべてある | 診断なし |
| `missing_left` | 境界値未満の入力がない | `ShippingFee: no test value less than 5000` |
| `missing_right` | 境界値より大きい入力がない | `ShippingFee: no test value greater than 5000` |
| `ignore_want` | `want`にある5000をテスト入力として数えない | `ShippingFee: boundary value 5000 is not tested` |
| `separate_functions` | テスト値を関数名ごとに分けて照合する | `ShippingFee`の境界値不足だけを診断 |
| `unsupported` | 対象外の構文を誤って診断しない | 診断なし |
| `named_constant_missing` | 名前付き定数を境界値として評価する | `ShippingFee: boundary value 5000 is not tested` |
| `named_constant_complete` | 名前付き定数でも3分類が揃えば診断しない | 診断なし |
| `constant_expression` | 定数式を計算後の境界値として評価する | `ShippingFee: boundary value 5000 is not tested` |
| `runtime_variable` | 実行時に変化しうる変数は対象にしない | 診断なし |

## analyzerを各ケースへ実行する

Step 3Cを完了した後、リポジトリのルートで実行します。診断があるケースでは、最後に表示される`exit status 3`も想定どおりです。

```bash
go run ./cmd/boundary ./boundary/testdata/missing_boundary
go run ./cmd/boundary ./boundary/testdata/complete
go run ./cmd/boundary ./boundary/testdata/missing_left
go run ./cmd/boundary ./boundary/testdata/missing_right
go run ./cmd/boundary ./boundary/testdata/ignore_want
go run ./cmd/boundary ./boundary/testdata/separate_functions
go run ./cmd/boundary ./boundary/testdata/unsupported
go run ./cmd/boundary ./boundary/testdata/named_constant_missing
go run ./cmd/boundary ./boundary/testdata/named_constant_complete
go run ./cmd/boundary ./boundary/testdata/constant_expression
go run ./cmd/boundary ./boundary/testdata/runtime_variable
```

すべてのケースを自動テストとしてまとめて確認する場合は、次を実行します。

```bash
go test -tags=workshop_solution ./boundary -run '^TestAnalyzer$' -v
```

1つのケースだけを自動テストで確認する場合は、サブテスト名を指定します。

```bash
go test -tags=workshop_solution ./boundary -run 'TestAnalyzer/missing_boundary' -v
```

`missing_boundary`の部分は、表にあるディレクトリ名へ置き換えられます。
