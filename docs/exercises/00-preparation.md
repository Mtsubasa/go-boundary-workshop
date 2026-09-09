# Ex0: 準備とpreflight

## この章のゴール

当日使うrepositoryを取得し、codeを変更する前のtestが成功することを確認します。

## 必要なもの

- Go codeを編集・実行でき、Gitを利用できるノートPC
- ノートPCの電源アダプター

GitHub account、ASTや静的解析の知識、特定のeditorは必要ありません。

## repositoryを取得する

```bash
git clone https://github.com/Mtsubasa/go-boundary-workshop.git
cd go-boundary-workshop
go mod download
go test ./...
```

すべてのpackageで`ok`または`[no test files]`と表示されれば準備完了です。

## うまくいかない場合

次の結果とerror全文を講師へ共有してください。

```bash
go version
git --version
go env GOMOD
```

当日は講師が復旧用ZIPを用意します。通常はbranch操作をせず、cloneしたworking treeのままExerciseを進めます。
