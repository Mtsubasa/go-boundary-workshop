# 04. Recovery Guide

Workshop中にcodeが動かなくなった場合は、作業を消さずにcheckpointから必要なfileだけ復元できます。branchを切り替える必要はありません。

## まず現在の作業を保存する

復元操作の前に、変更したfileを確認します。

```bash
git status --short
git diff
```

現在の作業をstashへ保存します。`-u`を付けることで、新しく作ったfileも保存対象になります。

```bash
git stash push -u -m "workshop: before recovery"
```

このあとに行う`git restore`は対象fileをcheckpointの内容へ置き換えます。stashを作成してから実行してください。

## Stepごとの復元command

すべてrepositoryのroot directoryで実行します。

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

### Step 3CのCore実装へ復元する

```bash
git restore --source=core-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

`ShippingFee: boundary value 5000 is not tested`と表示されれば復元できています。

`shipping.go`や`shipping_test.go`も含めて、Coreのdemo開始状態へ戻す場合は次を使います。

```bash
git restore --source=core-complete -- \
  boundary/analyzer.go \
  examples/shipping/shipping.go \
  examples/shipping/shipping_test.go
```

### Types Extra完了状態へ復元する

```bash
git restore --source=types-complete -- boundary/analyzer.go
go run ./cmd/boundary ./examples/shipping
```

named constantのexampleを含むcheckpoint全体を確認したい場合は、working treeを変更せずに次の差分を確認できます。

```bash
git diff core-complete..types-complete -- boundary/analyzer.go
```

## checkpointを確認するだけの場合

作業fileを置き換えず、完成状態との差分だけを表示できます。

```bash
git diff main..ex03-complete -- boundary/analyzer.go
git diff ex03-complete..ex04-complete -- boundary/analyzer.go
git diff ex04-complete..core-complete -- boundary/analyzer.go
git diff core-complete..types-complete -- boundary/analyzer.go
```

まず差分を見て、自分のcodeを直せそうなら`git restore`は不要です。

## stashへ保存した作業を戻す

保存した作業は次のcommandで確認できます。

```bash
git stash list
git stash show -p stash@{0}
```

checkpointのcodeを使わず、保存前の状態へ完全に戻りたい場合は、先にcheckpointから復元した変更を取り消し、その後でstashを戻します。

```bash
git restore --source=HEAD -- \
  boundary/analyzer.go \
  examples/shipping/shipping.go \
  examples/shipping/shipping_test.go
git stash pop
```

`git stash pop`でconflictが表示された場合、stashは削除されず残ります。慌ててfileを削除せず、講師へ表示されたmessageを共有してください。

## 復元後の確認

どのcheckpointを使った場合も、最後に次を実行します。

```bash
gofmt -w boundary/analyzer.go
go test ./...
git status --short
```

`go test ./...`が成功し、`git status --short`に自分が意図したfileだけ表示されていれば作業を再開できます。
