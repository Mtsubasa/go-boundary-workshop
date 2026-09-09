# CI Integration

完成した`analysis.Analyzer`は、local CLIだけでなく`go vet`のcustom toolとしてCIから実行できます。このページはWorkshop後に自分のrepositoryへ組み込むための手順です。当日のハンズオン中にworkflow fileを入力する必要はありません。

## 1. localで動作を確認する

まずanalyzerをbuildします。

```bash
go build -o boundary-checker ./cmd/boundary
```

buildしたbinaryを`go vet`へ渡します。

```bash
go vet -vettool="$(pwd)/boundary-checker" ./examples/shipping
```

境界値5000がtest tableにない場合は、次のような診断と0以外の終了statusになります。

```text
examples/shipping/shipping.go:4:13: ShippingFee: boundary value 5000 is not tested
```

`examples/shipping/shipping_test.go`へ次を追加すると、診断が消えてstatus 0になります。

```go
{input: 5000, want: 0},
```

## 2. GitHub Actions workflowを追加する

自分のrepositoryに`.github/workflows/boundary-check.yml`を作り、次を保存します。

```yaml
name: Boundary check

on:
  pull_request:
  push:
    branches:
      - main

permissions:
  contents: read

jobs:
  boundary-check:
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@v7

      - name: Set up Go
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache: true

      - name: Run tests
        run: go test ./...

      - name: Build boundary checker
        run: go build -o "${RUNNER_TEMP}/boundary-checker" ./cmd/boundary

      - name: Check boundary test cases
        run: go vet -vettool="${RUNNER_TEMP}/boundary-checker" ./...
```

`go-version-file: go.mod`により、repositoryの`go.mod`に書かれたGo versionを利用します。`permissions: contents: read`はcheckoutに必要な最小権限です。

## 3. CIの結果を読む

| 結果 | 意味 |
|---|---|
| `Run tests`が失敗 | 通常のtestが失敗している |
| `Check boundary test cases`が失敗 | analyzerが境界値caseの不足を診断した |
| 両方成功 | 通常testと、今回実装した静的ruleの両方を通過した |

静的解析はtestを実行しません。`go test`の成否と、境界値caseの有無は別の確認なので、workflowでもstepを分けています。

## 4. 自分のprojectへ持ち込むとき

Workshop repositoryとdirectory構成が違う場合は、次の2箇所を変更します。

- `./cmd/boundary`: analyzer CLIがあるpackage
- `./...`: 解析したいpackage pattern

たとえば`./internal/analyzers/cmd/projectlint`にCLIがあり、`./service/...`だけを検査する場合は次のようにします。

```yaml
- name: Build project analyzer
  run: go build -o "${RUNNER_TEMP}/projectlint" ./internal/analyzers/cmd/projectlint

- name: Run project analyzer
  run: go vet -vettool="${RUNNER_TEMP}/projectlint" ./service/...
```

## 当日の説明範囲

Workshopでは、localで成功した`go vet -vettool`とCIの2 stepが同じcommandであることだけを確認します。GitHub Actionsの文法や権限設計は主題ではないため、詳細はこのページを後から参照してください。
