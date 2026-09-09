# Analyzer test fixtures

このdirectoryには、完成したanalyzer自身をtestするための入力codeを置いています。

`examples/shipping`はWorkshop中に編集して動きを確認する題材です。一方、ここにある各directoryは編集せず、analyzerが期待どおりの診断を出すか確認するために使います。

## 今回確認するcase

| directory | 確認すること | 期待する診断 |
|---|---|---|
| `missing_boundary` | 境界の左右はあるが、境界値そのものがない | `ShippingFee: boundary value 5000 is not tested` |
| `complete` | less / boundary / greaterがすべてある | 診断なし |
| `missing_left` | 境界値未満の入力がない | `ShippingFee: no test value less than 5000` |
| `missing_right` | 境界値より大きい入力がない | `ShippingFee: no test value greater than 5000` |
| `ignore_want` | `want`にある5000をtest入力として数えない | `ShippingFee: boundary value 5000 is not tested` |
| `separate_functions` | test値を関数名ごとに分けて照合する | `ShippingFee`の境界値不足だけを診断 |
| `unsupported` | Coreの対象外にした構文を誤って診断しない | 診断なし |

## analyzerを各caseへ実行する

Step 3Cを完了した後、repositoryのroot directoryで実行します。診断があるcaseでは、最後に表示される`exit status 3`も想定どおりです。

```bash
go run ./cmd/boundary ./boundary/testdata/missing_boundary
go run ./cmd/boundary ./boundary/testdata/complete
go run ./cmd/boundary ./boundary/testdata/missing_left
go run ./cmd/boundary ./boundary/testdata/missing_right
go run ./cmd/boundary ./boundary/testdata/ignore_want
go run ./cmd/boundary ./boundary/testdata/separate_functions
go run ./cmd/boundary ./boundary/testdata/unsupported
```

すべてのcaseを自動testとしてまとめて確認する場合は、次を実行します。

```bash
go test -tags=workshop_solution ./boundary -run '^TestAnalyzer$' -v
```

1つのcaseだけを自動testで確認する場合は、subtest名を指定します。

```bash
go test -tags=workshop_solution ./boundary -run 'TestAnalyzer/missing_boundary' -v
```

`missing_boundary`の部分は、表にあるdirectory名へ置き換えられます。
