# go-backend-template

<a href="https://xcfile.dev"><img src="https://xcfile.dev/badge.svg" alt="xc compatible" /></a>

ハッカソンなど短期間でWebアプリを開発する際のバックエンドのGo実装例です。
学習コストと開発コストを抑えることを目的としています。

## How to use

GitHubの `Use this template` ボタンからレポジトリを作成します。

[`gonew`](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) コマンドからでも作成できます。`gonew` コマンドを使うと、モジュール名を予め変更した状態でプロジェクトを作成することができます。

```sh
gonew github.com/ras0q/go-backend-template {{ project_name }}
```

## Requirements

最低限[Docker](https://www.docker.com/)と[Docker Compose](https://docs.docker.com/compose/)が必要です。
[Compose Watch](https://docs.docker.com/compose/file-watch/)を使うため、Docker Composeのバージョンは2.22以上にしてください。

Linter, Formatterには[golangci-lint](https://golangci-lint.run/)を使っています。
VSCodeを使用する場合は`.vscode/settings.json`でLinterの設定を行ってください

```json
{
  "go.lintTool": "golangci-lint"
}
```

## Tasks

開発に用いるコマンド一覧

> [!TIP]
> `xc` を使うことでこれらのコマンドを簡単に実行できます。
> 詳細は以下のページをご覧ください。
>
> - [xc](https://xcfile.dev)
> - [MarkdownベースのGo製タスクランナー「xc」のススメ](https://zenn.dev/trap/articles/af32614c07214d)
>
> ```bash
> go install github.com/joerdav/xc/cmd/xc@latest
> ```

### Build

アプリをビルドします。

```sh
CMD=server
go mod download
go build -o ./bin/${CMD} ./cmd/${CMD}
```

### Dev

ホットリロードの開発環境を構築します。

```sh
docker compose watch
```

API、DB、DB管理画面が起動します。
各コンテナが起動したら、以下のURLにアクセスすることができます。
Compose Watchにより、ソースコードの変更を検知して自動で再起動します。

- <http://localhost:8080/> (API)
- <http://localhost:8081/> (DBの管理画面)

### Test

全てのテストを実行します。

```sh
go test -v -cover -race -shuffle=on ./...
```

### Test-Unit

単体テストを実行します。

```sh
go test -v -cover -race -shuffle=on ./internal/...
```

### Test-Integration

結合テストを実行します。

```sh
[ ! -e ./go.work ] && go work init . ./integration_tests
go test -v -cover -race -shuffle=on ./integration_tests/...
```

### Test-Integration:Update

結合テストのスナップショットを更新します。

```sh
[ ! -e ./go.work ] && go work init . ./integration_tests
go test -v -cover -race -shuffle=on ./integration_tests/... -update
```

### Lint

Linter (golangci-lint) を実行します。

```sh
golangci-lint run --timeout=5m --fix ./...
```

## Directory structure

[Organizing a Go module - The Go Programming Language](https://go.dev/doc/modules/layout#server-project) などを参考にしています。

```bash
$ tree -d
.
├── bin # ビルドしたバイナリ
├── cmd # エントリーポイント
│   └── server # サーバーのエントリーポイント (main.goを置く)
│       └── server # サーバー固有の設定
├── integration_tests # 結合テスト
├── internal # アプリケーション本体のロジック
│   ├── handler # ルーティング
│   └── repository # DBアクセス
└── pkg # 汎用パッケージ
    ├── config # アプリ・DBの設定
    └── database # DBの初期化、マイグレーション
        └── migrations # DBマイグレーションのスキーマ

12 directories
```

特に重要なものは以下の通りです。

### `cmd/`

アプリケーションのエントリーポイントを配置します。
エントリーポイントは複数のアプリケーションを持つことも可能です。

テンプレートではサーバーのエントリーポイントが`cmd/server/main.go`に配置されています。
`cmd/server/server/` にはDIやヘルスチェックなど、サーバー固有の設定を書いています。

### `internal/`

アプリケーション本体のロジックを配置します。
主に2つのパッケージに分かれています。

- `handler/`: ルーティング
  - 飛んできたリクエストを裁いてレスポンスを生成する
  - DBアクセスは`repository/`で実装したメソッドを呼び出す
  - **Tips**: リクエストのバリデーションがしたい場合は↓のどちらかを使うと良い
    - [go-playground/validator](https://github.com/go-playground/validator): タグベースのバリデーション
    - [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation): コードベースのバリデーション
- `repository/`: ストレージ操作
  - DBや外部ストレージなどのストレージにアクセスする
    - 引数のバリデーションは`handler/`に任せる

**Tips**: `internal`パッケージは他モジュールから参照されません（参考: [Go 1.4 Release Notes](https://go.dev/doc/go1.4#internalpackages)）。
依存性注入や外部ライブラリの初期化のみを`cmd/`や`pkg`で公開し、アプリケーションのロジックは`internal/`に閉じることで、後述の`integration_tests/go.mod`などの外部モジュールからの参照を最小限にすることができ、開発の効率を上げることができます。

### `pkg/`

汎用的なパッケージを配置します。

- `config/`: アプリ・DBの設定
  - 環境変数を読み込むためのパッケージ
- `database/`: DBの初期化、マイグレーション
  - DBのスキーマを定義する
  - **Tips**: マイグレーションツールは[pressly/goose](https://github.com/pressly/goose)を使っている

### `integration_tests/`

結合テストを配置します。
APIエンドポイントに対してリクエストを送り、レスポンスを検証します。
短期開発段階では時間があれば書く程度で良いですが、長期開発に向けては書いておくと良いでしょう。

```go
package integration_tests

import (
  "testing"
  "gotest.tools/v3/assert"
)

func TestUser(t *testing.T) {
  t.Run("get users", func(t *testing.T) {
    t.Run("success", func(t *testing.T) {
      t.Parallel()
      rec := doRequest(t, "GET", "/api/v1/users", "")

      expectedStatus := `200 OK`
      expectedBody := `[{"id":"[UUID]","name":"test","email":"test@example.com"}]`
      assert.Equal(t, rec.Result().Status, expectedStatus)
      assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
    })
  })
}
```

**Tips**: DBコンテナの立ち上げには[ory/dockertest](https://github.com/ory/dockertest)を使っています。

**Tips**: アサーションには[gotest.tools](https://github.com/gotestyourself/gotest.tools)を使っています。
`go test -update`を実行することで、`expectedXXX`のスナップショットを更新することができます（参考: [gotest.toolsを使う - 詩と創作・思索のひろば](https://motemen.hatenablog.com/entry/2022/03/gotest-tools)）。

外部サービス（traQ, Twitterなど）へのアクセスが発生する場合はTest Doublesでアクセスを置き換えると良いでしょう。

## Improvements

長期開発に向けた改善点をいくつか挙げておきます。

- ドメインを書く (`internal/domain/`など)
  - 現在は簡単のためにAPIスキーマとDBスキーマのみを書きこれらを直接やり取りしている
  - 本来はアプリの仕様や概念をドメインとして書き、スキーマの変換にはドメインを経由させるべき
- クライアントAPIスキーマを共通化させる
  - OpenAPIやGraphQLを使い、そこからGoのファイルを生成する
- 単体テスト・結合テストのカバレッジを上げる
  - カバレッジの可視化には[Codecov](https://codecov.io)(traPだと主流)や[Coveralls](https://coveralls.io)が便利
- ログの出力を整備する
  - ロギングライブラリは好みに合ったものを使うと良い
