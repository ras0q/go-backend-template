# frontend

Webフロントエンドのサンプルとして [Vite](https://vite.dev) で作成したテンプレートを置いています。

他のframeworkを使用したい場合、このディレクトリを削除しフロントエンドプロジェクトを作成しなおしてください。

`static.go` はGoのバックエンドにフロントエンドの成果物を埋め込むために必要なので消さないでください。

```sh
# in project root

cp ./frontend/static.go /tmp/static.go
rm -r ./frontend

# See https://vite.dev/guide/
npm create vite@latest

cp /tmp/static.go ./frontend/static.go
```
