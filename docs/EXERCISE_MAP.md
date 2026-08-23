# Exercise Map

Core完成版を実装して実測した後、このmapをもとにstarterとcheckpointを切り出します。時間は仮置きです。

| 区分 | Exercise | 参加者が編集する中心 | 完了確認 |
|---|---|---|---|
| Core | Ex1: AST Viewer | 編集なし | 4種類のnodeを画面で見つける |
| Core | Ex2: Analyzer起動 | Analyzer metadataとrun関数 | commandがerrorなく終了する |
| Core | Ex3: 関数を報告 | `FuncDecl`の走査と`Reportf` | `found function ShippingFee` |
| Core | Ex4: 境界値を抽出 | `IfStmt`、`BinaryExpr`、`BasicLit` | `found boundary value 5000` |
| Core | Ex5: test値を収集 | `TestXxx`と`input`の`KeyValueExpr` | 4999と5001を取得する |
| Core | Ex6: 照合して警告 | left/boundary/rightの判定 | boundary不足を1件報告する |
| Core | Ex7: analyzerをtest | fixtureと期待診断 | 不足あり・なしがtestで固定される |
| Extra | Extra 1: named constant | `TypesInfo.Types`と`go/constant` | `const Boundary`を評価する |
| Extra | Extra 2: 複合条件 | 式treeの再帰 | `&&`、`||`内の比較を取得する |
| Extra | Extra 3: operator | 比較の正規化 | 4つの大小比較を共通化する |

## Core完成版からの切り出し方

| checkpoint | 完成している範囲 | 次に参加者が触る箇所 |
|---|---|---|
| `main` | CLI、Analyzer metadata、題材code、test harness | `FuncDecl`の走査 |
| `ex03-complete` | 関数名の一時診断 | `IfStmt`以下の走査 |
| `ex04-complete` | 境界値5000の一時診断 | test値の収集と3分類 |
| `core-complete` | Coreの全診断と自動test | Extraまたは教材の読み進め |
| `types-complete` | named constantとconstant expression | 一般変数に必要な値追跡を考える |

Core完成版で検証用helperは実装済みだが、すべてを参加者に書かせない。特にfixture読込、整数literalのparse、対応外構文のfilterはstarterで提供する候補とする。

2026-08-23時点で上記5 branchは作成・実測済み。通常参加者は`main`だけを使い、checkpoint branchは講師の復旧用とする。

## Checkpoint方針

- 参加者はmainのstarterをcloneし、同じworking treeで順番に編集する
- Exerciseごとのbranch切り替えは要求しない
- 講師は`ex03-complete`、`ex04-complete`、`core-complete`のsnapshotを持つ
- 復旧用snapshotはRelease zipまたは明示的なcopy手順で提供する
- 各Exerciseには「期待出力」「hint 1〜3」「復旧先」を記載する

## 分割時の判断基準

- 1 Exerciseで新しく登場する主要AST nodeは1〜2種類まで
- 参加者が記述するcodeは一度に10〜20行程度まで
- 10分を超えるboilerplateはstarter側へ置く
- file探索、package variant、error handlingなど学習目標でない複雑さは提供済みにする
- 途中の一時的な診断は、次のStepで置き換えてよい
