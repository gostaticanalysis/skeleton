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

#### 解析器

`x/tools/go/analysis`パッケージを用いて開発された静的解析ツールは、`*analysis.Analyzer`型の値として表現されます。mylinterの場合、`mylinter.go`に`Analyzer`変数として定義されています。

生成されたコードは、`inspect.Analyzer`を用いた簡単な静的解析ツールを実装しています。この静的解析ツールは、`gopher`という名前の識別子を見つけるだけです。

```go
package mylinter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "mylinter is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "mylinter",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pass.Reportf(n.Pos(), "identifier is gopher")
			}
		}
	})

	return nil, nil
}
```

#### テストコード

skeletonは、テストコードも生成します。`x/tools/go/analysis`パッケージはサブパッケージの`analysistest`パッケージとして、テストライブラリを提供しています。`analysistest.Run`関数は`testdata/src`ディレクトリ以下にあるソースコードを使ってテストを実行します。この関数の第2引数はテストデータのディレクトリです。第3引数はテスト対象のAnalyzer、第4引数以降はテストデータとして利用するパッケージ名です。

```go
package mylinter_test

import (
	"testing"

	"github.com/gostaticanalysis/example.com/mylinter"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, mylinter.Analyzer, "a")
}
```

mylinterの場合、テストは`testdata/src/a/a.go`ファイルをテストデータとして利用します。`mylinter.Analyzer`は`gopher`識別子をソースコードの中から探し報告します。テストでは、期待する報告をコメントで記述します。コメントは`want`で始まり、その後に期待するメッセージが正規表現で記述されます。テストは期待するメッセージで報告がされなかったり、期待していない報告がされた場合に失敗します。

```go
package a

func f() {
	// The pattern can be written in regular expression.
	var gopher int // want "pattern"
	print(gopher)  // want "identifier is gopher"
}
```

デフォルトでは`go mod tidy`コマンドと`go test`コマンドを実行すると、テストは失敗します。これは`pattern`というメッセージで作った静的解析ツールが報告をしないためです。

```
$ go mod tidy
go: finding module for package golang.org/x/tools/go/analysis
go: finding module for package github.com/gostaticanalysis/testutil
go: finding module for package golang.org/x/tools/go/analysis/passes/inspect
go: finding module for package golang.org/x/tools/go/analysis/unitchecker
go: finding module for package golang.org/x/tools/go/ast/inspector
go: finding module for package golang.org/x/tools/go/analysis/analysistest
go: found golang.org/x/tools/go/analysis in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/analysis/passes/inspect in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/ast/inspector in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/analysis/unitchecker in golang.org/x/tools v0.1.10
go: found github.com/gostaticanalysis/testutil in github.com/gostaticanalysis/testutil v0.4.0
go: found golang.org/x/tools/go/analysis/analysistest in golang.org/x/tools v0.1.10

$ go test
--- FAIL: TestAnalyzer (0.06s)
    analysistest.go:454: a/a.go:5:6: diagnostic "identifier is gopher" does not match pattern `pattern`
    analysistest.go:518: a/a.go:5: no diagnostic was reported matching `pattern`
FAIL
exit status 1
FAIL	github.com/gostaticanalysis/example.com/mylinter	1.270s
```

#### 実行可能ファイル

skeletonは`cmd`ディレクトリ以下に`main.go`も生成します。この`main.go`をビルドし生成した実行可能ファイルは、`go vet`コマンド経由で実行される必要があります。`go vet`コマンドの`-vettool`オプションは生成した実行可能ファイルへの絶対パスを指定します。

```
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

`singlechecker`パッケージは、単一のAnalyzerを実行するためのパッケージで`go vet`は必要としません。利用するには`-checker=single`を指定します。

`multichecker`パッケージは、複数のAnalyzerを実行するためのパッケージで`go vet`は必要としません。利用するには`-checker=multi`を指定します。

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
