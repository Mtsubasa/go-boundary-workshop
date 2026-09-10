# 03. 復旧手順

ワークショップ中にコードが動かなくなった場合は、作業を消さずにチェックポイントから必要なファイルだけ復元できます。ブランチを切り替える必要はありません。

## まず現在の作業を保存する

復元操作の前に、変更したファイルを確認します。

```bash
git status --short
git diff
```

現在の作業をstashへ保存します。`-u`を付けることで、新しく作ったファイルも保存対象になります。

```bash
git stash push -u -m "workshop: before recovery"
```

このあとに行う`git restore`は、対象ファイルをチェックポイントの内容へ置き換えます。stashを作成してから実行してください。

## Stepごとの復元コマンド

すべてリポジトリのルートで実行します。

### Step 1完了状態へ復元する

```bash
git restore --source=ex03-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

`found function ShippingFee`と表示されれば復元できています。

### Step 2完了状態へ復元する

```bash
git restore --source=ex04-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

`ShippingFee: found boundary value 5000`と表示されれば復元できています。

### Step 3Cの完成状態へ復元する

```bash
git restore --source=step3-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

`ShippingFee: boundary value 5000 is not tested`と表示されれば復元できています。

`shipping.go`や`shipping_test.go`も含めて、Step 3Cの確認に必要な状態へ戻す場合は次を使います。

```bash
git restore --source=step3-complete -- \
  boundary/analyzer.go \
  examples/shipping/shipping.go \
  examples/shipping/shipping_test.go
```

### 応用課題3の完了状態へ復元する

```bash
git restore --source=types-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

名前付き定数を含むチェックポイント全体を確認したい場合は、作業中のファイルを変更せずに次の差分を確認できます。

```bash
git diff step3-complete..types-complete -- boundary/analyzer.go
```

## チェックポイントを確認するだけの場合

作業中のファイルを置き換えず、完成状態との差分だけを表示できます。

```bash
git diff main..ex03-complete -- boundary/analyzer.go
git diff ex03-complete..ex04-complete -- boundary/analyzer.go
git diff ex04-complete..step3-complete -- boundary/analyzer.go
git diff step3-complete..types-complete -- boundary/analyzer.go
```

まず差分を見て、自分のコードを直せそうなら`git restore`は不要です。

## stashへ保存した作業を戻す

保存した作業は次のコマンドで確認できます。

```bash
git stash list
git stash show -p stash@{0}
```

チェックポイントのコードを使わず、保存前の状態へ完全に戻りたい場合は、先にチェックポイントから復元した変更を取り消し、その後でstashを戻します。

```bash
git restore --source=HEAD -- \
  boundary/analyzer.go \
  examples/shipping/shipping.go \
  examples/shipping/shipping_test.go
git stash pop
```

`git stash pop`で競合が表示された場合、stashは削除されず残ります。慌ててファイルを削除せず、表示された内容を講師へ共有してください。

## 復元後の確認

どのチェックポイントを使った場合も、最後に次を実行します。

```bash
gofmt -w boundary/analyzer.go
go test ./...
git status --short
```

`go test ./...`が成功し、`git status --short`に自分が意図したファイルだけ表示されていれば作業を再開できます。
