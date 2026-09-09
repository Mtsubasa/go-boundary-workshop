# 02. CI Integration

この内容はCore完了後のOptionalです。実装した`analysis.Analyzer`は、CLIとしての直接実行に加え、`go vet`のcustom toolとしてCIから実行できます。以下の例は`core-complete` branch相当の実装を前提とします。

Workshop当日に時間が限られる場合は、「localで動作を確認する」までをdemoし、GitHub Actionsはworkflowの構造を説明するだけで十分です。workflowの作成、push、実行待ちは必須作業に含めません。

## localで動作を確認する

analyzerをbuildします。

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

## GitHub Actionsで実行する

自分のrepositoryへ導入する場合は、`.github/workflows/boundary-check.yml`を作成します。

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

`go test`と静的解析は異なる観点を確認するため、別々のstepとして実行します。

`go vet -vettool`は、診断があると非0で終了します。その終了statusをGitHub Actionsが受け取り、`Check boundary test cases` stepを失敗として表示します。

## CIの結果を読む

| 結果 | 意味 |
|---|---|
| `Run tests`が失敗 | 通常のtestが失敗している |
| `Check boundary test cases`が失敗 | analyzerが境界値caseの不足を診断した |
| 両方成功 | 通常testと静的ruleの両方を通過した |

## 別のrepositoryへ組み込む

directory構成に合わせて、analyzer CLIのpackageと解析対象のpackage patternを変更します。

```yaml
- name: Build project analyzer
  run: go build -o "${RUNNER_TEMP}/projectlint" ./internal/analyzers/cmd/projectlint

- name: Run project analyzer
  run: go vet -vettool="${RUNNER_TEMP}/projectlint" ./service/...
```
