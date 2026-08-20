# Ex2: analyzerを起動する

## この章のゴール

まだ何も診断しない最小の`analysis.Analyzer`をCLIとして実行します。

## 3つの登場人物

| 名前 | 今回の役割 |
|---|---|
| `singlechecker` | `Analyzer`をcommandとして起動する |
| `Analyzer` | toolの名前、説明、実行する関数を定義する |
| `Pass` | ASTや型情報を受け取り、診断を報告する窓口になる |

詳しいAPIをすべて理解してから始める必要はありません。最初は「CLIが`run`を呼び、`pass.Reportf`で診断を返せる」と捉えます。

## 実行する

```bash
go run ./cmd/boundary ./examples/shipping
```

starterではまだ診断がなく、errorなしで終了すれば成功です。

## codeをつなげて読む

1. `cmd/boundary/main.go`の`singlechecker.Main`
2. `boundary/analyzer.go`の`Analyzer`
3. `Analyzer.Run`に設定された`run`

## 完了条件

- commandを実行できた
- `main -> Analyzer -> run`の順につながっていることを確認した

