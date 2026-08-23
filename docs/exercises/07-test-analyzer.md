# Ex7: analyzer自身をtestする

## この章のゴール

手作業で試した1例だけでなく、対応範囲の重要なcaseをfixtureで固定します。

test harnessはstarterに用意します。当日はtest codeを一から書かず、何を固定しているかを読んで実行します。preflight時はまだanalyzerが未完成なのでbuild tagで除外しています。

## 実行する

```bash
go test -tags workshop_solution ./boundary -v
```

## 固定するcase

- boundaryだけ不足する
- lessだけ不足する
- greaterだけ不足する
- 3分類がすべて揃う
- `want: 5000`をinputとして扱わない
- 別functionのtest値を混ぜない
- Core非対応の構文を対象にしない

## fixtureを見る

`boundary/testdata/missing_boundary`と`boundary/testdata/complete`を比較します。違いは`input: 5000`の1行です。

## 完了条件

- `go test -tags workshop_solution ./boundary -v`が成功する
- test caseの追加で診断が消えることを自動testでも確認できた

Coreはここで完成です。全員がExtraへ進む必要はありません。
