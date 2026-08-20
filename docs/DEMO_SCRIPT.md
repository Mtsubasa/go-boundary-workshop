# Complete Demo Script

80〜90分の「完全版demo・CI組み込み・まとめ」で使う講師用手順です。説明しながら操作しても5〜6分に収め、残りをまとめと質疑に使います。

## 0. 事前確認

登壇前にrepository rootで次を実行します。

```bash
go test ./...
go build -o boundary-checker ./cmd/boundary
```

terminalはrepository rootで開き、文字sizeを大きくします。`examples/shipping/shipping_test.go`は5000がない初期状態に戻しておきます。

## 1. coverageは100%になる

```bash
go test -cover ./examples/shipping
```

見せる出力:

```text
coverage: 100.0% of statements
```

話す要点:

- 4999で`if`の中、5001で`if`の外を実行している
- statement coverageの観点では両方を通っている
- しかし5000で条件が切り替わること自体はtestしていない

## 2. 自作analyzerは不足を報告する

```bash
go run ./cmd/boundary ./examples/shipping
```

見せる出力:

```text
examples/shipping/shipping.go:4:13: ShippingFee: boundary value 5000 is not tested
exit status 3
```

話す要点:

- 一般的なcoverageとは別の、自分が守りたい観点を検査している
- 診断位置はASTで見つけた`5000`の位置
- 終了statusが0以外なので、人の目だけでなくCIでも止められる

## 3. test caseを1つ追加する

`examples/shipping/shipping_test.go`へ次を追加します。

```go
{input: 5000, want: 0},
```

もう一度実行します。

```bash
go test ./examples/shipping
go run ./cmd/boundary ./examples/shipping
```

testは成功し、analyzerの診断も消えます。

話す要点:

- source codeを直したのではなく、toolが求めたtestを追加した
- 自分で書いたASTの処理が、具体的なfeedbackになった
- 今日覚えてほしい面白さはこの瞬間

## 4. `go vet`へ組み込む

```bash
go build -o boundary-checker ./cmd/boundary
go vet -vettool="$(pwd)/boundary-checker" ./examples/shipping
```

話す要点:

- `analysis.Analyzer`として作ると、単体CLIだけでなく`go vet`の仕組みに載せられる
- 当日はCI設定の文法を詳しく説明しない

## 5. CIの完成イメージ

slideまたは教材では、次の最小例だけを見せます。

```yaml
- name: Build boundary checker
  run: go build -o boundary-checker ./cmd/boundary

- name: Check boundary test cases
  run: go vet -vettool="$(pwd)/boundary-checker" ./...
```

Workshop中にworkflow fileを最初から入力させません。参加者には「同じcommandをCIで実行できる」ことだけ伝えます。

## 6. まとめの言葉

次の3点で締めます。

1. ASTはsource codeをnodeとしてたどれる形にしたもの
2. `analysis`を使うと、見つけたnodeをsource位置付きの診断にできる
3. 既製toolが守らない観点でも、自分のこだわりを機械的な検査にできる

最後の一文:

> ASTのnode名を今日覚える必要はありません。次に「これを機械的に守れないかな」と思ったとき、AST Viewerと今日のrepositoryを思い出してもらえたら十分です。

## 当日のfallback

- network不調: 事前buildした`boundary-checker`とrepositoryのRelease ZIPを使う
- 参加者の実装が未完成: `core-complete`のsnapshotでdemoを続ける
- live editに失敗: 5000を含む完成済みtest fixtureを別fileで用意して切り替える
- 時間超過: `go vet`の実行を省き、CI snippetをslideで見せる
