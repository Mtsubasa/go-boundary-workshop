# 02. 応用課題5：CIへの組み込み

この内容はStep 3C完了後の応用課題です。実装した`analysis.Analyzer`は、コマンドとしての直接実行に加え、`go vet`の追加ツールとしてCIから実行できます。

GitHubへpushできるリポジトリがない場合は、「手元で動作を確認する」ところまでで完了として構いません。

## 手元で動作を確認する

analyzerをビルドします。

```bash
go build -o boundary-checker ./cmd/boundary
```

ビルドした実行ファイルを`go vet`へ渡します。

```bash
go vet -vettool="$(pwd)/boundary-checker" ./examples/shipping
```

Step 3Cで境界値5000を追加した状態では、何も表示せず終了ステータス0になります。

失敗も確認する場合は、`examples/shipping/shipping_test.go`の5000を一時的にコメントアウトして同じコマンドを実行します。次のような診断を表示し、終了ステータスは0以外になります。

```text
examples/shipping/shipping.go:4:13: ShippingFee: boundary value 5000 is not tested
```

確認後は5000のケースを元に戻します。

## GitHub Actionsで実行する

自分のリポジトリへ導入する場合は、用意された設定を`.github/workflows`へコピーします。

```bash
mkdir -p .github/workflows
cp ci/boundary-check.yml .github/workflows/boundary-check.yml
```

コピーされる内容は次のとおりです。

```yaml
name: Boundary check

on:
  push:
  pull_request:
  workflow_dispatch:

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
        run: go vet -vettool="${RUNNER_TEMP}/boundary-checker" ./examples/shipping
```

`go test`と静的解析は異なる観点を確認するため、別々のStepとして実行します。解析対象を`./examples/shipping`へ限定しているのは、`boundary/testdata`に意図的な不足ケースが含まれているためです。

`go vet -vettool`は、診断があると0以外で終了します。その終了ステータスをGitHub Actionsが受け取り、`Check boundary test cases`を失敗として表示します。

## CIの結果を読む

| 結果 | 意味 |
|---|---|
| `Run tests`が失敗 | 通常のテストが失敗している |
| `Check boundary test cases`が失敗 | analyzerが境界値ケースの不足を診断した |
| 両方成功 | 通常のテストと静的解析の両方を通過した |

## 別のリポジトリへ組み込む

ディレクトリ構成に合わせて、analyzerのパッケージと解析対象のパッケージパターンを変更します。

```yaml
- name: Build project analyzer
  run: go build -o "${RUNNER_TEMP}/projectlint" ./internal/analyzers/cmd/projectlint

- name: Run project analyzer
  run: go vet -vettool="${RUNNER_TEMP}/projectlint" ./service/...
```
