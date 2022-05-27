[English Version](./README.md)

# skeleton 

skeletonはGoの静的解析ツールのためのスケルトンコードジェネレータです。[x/tools/go/analysis](https://golang.org/x/tools/go/analysis)パッケージや[x/tools/go/packages](https://golang.org/x/tools/go/packages)パッケージを用いた静的解析ツールの開発を簡単にします。

## x/tools/go/analysisパッケージ

[x/tools/go/analysis](https://golang.org/x/tools/go/analysis)パッケージは静的解析ツールをモジュール化するためのパッケージです。[analysis.Analyzer](https://golang.org/x/tools/go/analysis)型を1つの単位として扱います。

`x/tools/go/analysis`パッケージは、静的解析ツールの共通部分を定型化しています。skeletonは定型化されているコードの大部分をスケルトンコードとして生成します。`skeleton mylinter`コマンドを実行するだけで`*analyzer.Analyzer`型の初期化コードやテストコード、`go vet`から実行できる実行可能ファイルを作るための`main.go`を生成してくれます。

skeletonについて詳しく知りたい場合は、次のブログも参考になります。

* [skeletonで始めるGoの静的解析](https://engineering.mercari.com/blog/entry/20220406-eea588f493/)

`x/tools/go/analysis`パッケージやGoの静的解析自体を知りたい場合は、次の資料が参考になります。

* [プログラミング言語Go完全入門 14章 静的解析とコード生成](http://tenn.in/analysis)

## インストール

Goのバージョンによってインストール方法が異なります。

### Go1.16未満

```
$ go get -u github.com/gostaticanalysis/skeleton/v2
```

### Go1.16以上

```
$ go install github.com/gostaticanalysis/skeleton/v2@latest
```

## 使用方法

### モジュールパスを指定して作成

skeletonの引数にモジュールパスを指定するとそのパスでモジュールを生成します。ディレクトリ名はモジュールパスの最後の要素になります。`example.com/mylinter`と指定すると次のようになります。

```
$ skeleton example.com/mylinter
mylinter
├── cmd
│   └── mylinter
│       └── main.go
├── go.mod
├── mylinter.go
├── mylinter_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

`mylinter.go`に静的解析ツールの本体である`*analysis.Analyzer`型の変数が生成されます。`mylinter_test.go`はテストコードですが、基本的には変更する必要はありません。テストは`testdata`ディレクトリ以下に置かれたテスト対象のソースコードを用いて行われます。これは`x/tools/go/analysis`パッケージの仕様に基づいて行われます。詳しくは[ドキュメント](https://pkg.go.dev/golang.org/x/tools/go/analysis)または[プログラミング言語Go完全入門 14章 静的解析とコード生成](http://tenn.in/analysis)をご覧ください。

`cmd`ディレクトリ以下には`go vet`から実行される前提の実行可能ファイルを生成するための`main.go`が配置されています。生成した実行可能ファイルファイルは次のように`go vet`経由で実行します。引数は`go vet`と同じでパッケージなどを指定します。

```go
$ go vet -vettool=`which mylinter` ./...
```

### ディレクトリの上書き

すでにディレクトリが存在する場合は上書きするか聞かれます。

```
$ skeleton example.com/mylinter
mylinter is already exist, overwrite?
[1] No (Exit)
[2] Remove and create new directory
[3] Overwrite existing files with confirmation
[4] Create new files only
```

選んだ選択肢によって処理が変わります。

* [1] 上書きしない（終了）
* [2] 削除して新しいディレクトリを作成
* [3] すでにあるファイルを上書きするか都度確認
* [4] 新しいファイルのみ生成する

### cmdディレクトリを生成しない

`-cmd`オプションを`false`にすると`cmd`ディレクトリは生成されません。

```
$ skeleton -cmd=false example.com/mylinter
mylinter
├── go.mod
├── mylinter.go
├── mylinter_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

### go.modファイルを生成しない

skeletonはデフォルトでは`go.mod`ファイルを生成します。すでにGo Modules管理下にあるディレクトリでスケルトンコードを生成したい場合は、次のように`-gomod`オプションに`false`を指定します。

```
$ skeleton -gomod=false example.com/mylinter
mylinter
├── cmd
│   └── mylinter
│       └── main.go
├── mylinter.go
├── mylinter_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

### SKELETON_PREFIX環境変数

次のように`SKELETON_PREFIX`環境変数を指定するとモジュールパスの前にプリフィックスを付与します。

```
$ SKELETON_PREFIX=example.com skeleton mylinter
$ head -1 mylinter/go.mod
module example.com/mylinter
```

次のように[direnv](https://github.com/direnv/direnv)などを用いて特定のディレクトリ以下でプリフィックスをつけるようにすると便利です。

```
$ cat ~/repos/gostaticanalysis/.envrc
export SKELETON_PREFIX=github.com/gostaticanalysis
```

`SKELETON_PREFIX`環境変数を指定していても、`-gomod`オプションを`false`にした場合は親のモジュールのモジュールパスが使用されます。

### singlecheckerまたはmulticheckerの使用

デフォルトでは`main.go`では`go vet`から実行することを前提とした`unitchecker`パッケージが使われています。`-checker`オプションを指定することで、`singlechecker`パッケージや`multichecker`パッケージに変更できます。

`singlechecker`パッケージは、単一のAnalyzerを実行するためのパッケージで`go vet`は必要としません。利用するには`-change=single`を指定します。

`multichecker`パッケージは、複数のAnalyzerを実行するためのパッケージで`go vet`は必要としません。利用するには`-change=multi`を指定します。

次に`singlechecker`パッケージを利用した例を示します。

```
$ skeleton -checker=single example.com/mylinter
$ cat cmd/mylinter/main.go
package main

import (
		"mylinter"
		"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(mylinter.Analyzer) }
```

`singlechecker`パッケージや`multichecker`パッケージを利用した方が簡単そうに見えますが、`go vet`を使った恩恵を受けられないため、特にこだわりがない場合は`unitchecker`（デフォルト）を使用すると良いでしょう。

### スケルトンコードの種類を変更

skeletonは`-kind`オプションを指定することで生成するスケルトンコードを変更できます。

* `-kind=inspect`(デフォルト): `inspect.Analyzer`を用いたコードを作成
* `-kind=ssa`: `buildssa.Analyzer`で生成した静的単一代入(SSA, Static Single Assignment)形式を用いたコードを作成
* `-kind=codegen`: コード生成器を作成
* `-kind=packages`: `x/tools/go/packages`パッケージを用いたコードを作成

### コード生成器の作成

skeletonは`-kind`オプションに`codegen`を指定すると[gostaticanalysis/codegen](https://pkg.go.dev/github.com/gostaticanalysis/codegen)パッケージを用いたコード生成器のスケルトンコードも生成できます。

```
$ skeleton -kind=codegen example.com/mycodegen
mycodegen
├── cmd
│   └── mycodegen
│       └── main.go
├── go.mod
├── mycodegen.go
├── mycodegen_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            ├── go.mod
            └── mycodegen.golden
```

`gostaticanalysis/codegen`パッケージは実験的なパッケージです。ご注意ください。

### golangci-lintのプラグインを生成する

skeletonは`-plugin`パッケージを指定すると[golangci-lint](https://github.com/golangci/golangci-lint)からプラグインとして利用できるコードを生成します。

```
$ skeleton -plugin example.com/mylinter
mylinter
├── cmd
│   └── mylinter
│       └── main.go
├── go.mod
├── mylinter.go
├── mylinter_test.go
├── plugin
│   └── main.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

ビルド方法は[golangci-lintのドキュメント](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint)にも記載がありますが、生成されたコードの先頭にコメントとして記述されています。

```
$ skeleton -plugin example.com/mylinter
$ go build -buildmode=plugin -o path_to_plugin_dir example.com/mylinter/plugin/mylinter
```

もし、プラグインで特定のフラグを指定したい場合は、ビルドする際に`-ldflags`オプションを指定して設定します。この機能はskeletonで生成したコードのみに提供されます。詳しくは生成されたスケルトンコードをご覧ください。

```
$ skeleton -plugin example.com/mylinter
$ go build -buildmode=plugin -ldflags "-X 'main.flags=-funcs log.Fatal'" -o path_to_plugin_dir example.com/mylinter/plugin/mylinter
```

なお、プラグインは標準の`plugin`パッケージを使用するため、golangci-lintを`CGO_ENABLED=1`でビルドし直す必要があります。また、golangci-lintと生成した静的解析ツールで使用しているモジュールのバージョンを揃えないといけないため、あまりおすすめはしません。
