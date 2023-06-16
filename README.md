# go-backend-sample

ハッカソンなど短期間でWebアプリを開発する際のバックエンドのGo実装例です。
学習コストと開発コストを抑えることを目的としています。

## 使い方

- 最低限[Docker](https://www.docker.com/) ([Docker Compose](https://docs.docker.com/compose/))が必要です。
- linter, formatterには[golangci-lint](https://golangci-lint.run/)を使っています。
- 開発環境では[cosmtrek/air](https://github.com/cosmtrek/air)を使ったホットリロード開発が可能です
- makeコマンドのターゲット一覧とその説明は`make help`で確認できます

### 開発環境の実行

```sh
make dev
```

API、DB、DB管理画面が起動します。
各コンテナが起動したら、以下のURLにアクセスすることができます

- <http://localhost:8080/> (API)
- <http://localhost:8081/> (DBの管理画面)

### テストの実行

全てのテスト

```sh
make test
```

単体テストのみ

```sh
make test-unit
```

結合テストのみ

```sh
make test-integration
```

## 構成

- `main.go`: エントリーポイント
  - 依存ライブラリの初期化など最低限の処理のみを書く
  - ルーティングの設定は`./internal/handler/handler.go`に書く
  - 肥大化しそうなら`./internal/infrastructure/{pkgname}`を作って外部ライブラリの初期化処理を書くのもアリ
- `internal/`: アプリ本体の主実装
  - Tips: Goの仕様で`internal`パッケージは他プロジェクトから参照できない (<https://go.dev/doc/go1.4#internalpackages>)
  - `handler/`: ルーティング
    - 飛んできたリクエストを裁いてレスポンスを生成する
    - DBアクセスは`repository/`で実装したメソッドを呼び出す
    - Tips: リクエストのバリデーションがしたい場合は↓のどちらかを使うと良い
      - [go-playground/validator](https://github.com/go-playground/validator)でタグベースのバリデーションをする
      - [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation)でコードベースのバリデーションをする
  - `repository/`: DBアクセス
    - DBスキーマの定義とDBへのアクセス処理
      - 引数のバリデーションは`handler/`に任せる
    - 初期化スキーマは`schema.sql`に記述する
      - Tips: Goでは1.16から[embed](https://pkg.go.dev/embed)パッケージを使ってバイナリにファイルを文字列として埋め込むことができる
  - `pkg/`: 汎用パッケージ
    - 複数パッケージから使いまわせるようにする
    - 例: `pkg/config/`: アプリ・DBの設定
    - Tips: 外部にパッケージを公開したい場合は`internal/`の外に出しても良い
- `integration/`: 結合テスト
  - `internal/`の実装から実際にデータが取得できるかテストする
  - DBの立ち上げには[ory/dockertest](https://github.com/ory/dockertest)を使っている
  - 短期開発段階では時間があれば書く程度で良い
  - Tips: 外部サービス(traQ, Twitterなど)へのアクセスが発生する場合は[golang/mock](https://github.com/golang/mock)などを使ってモック(テスト用処理)を作ると良い

## 長期開発に向けた改善点

- ドメインを書く (`internal/domain/`など)
  - 現在は簡単のためにAPIスキーマとDBスキーマのみを書きこれらを直接やり取りしている
  - 本来はアプリの仕様や概念をドメインとして書き、スキーマの変換にはドメインを経由させるべき
- 単体テスト・結合テストのカバレッジを上げる
  - カバレッジの可視化には[Codecov](https://codecov.io)(traPだと主流)や[Coveralls](https://coveralls.io)が便利
- ログの出力を整備する
  - ロギングライブラリは好みに合ったものを使うと良い
