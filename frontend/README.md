# frontend

Webフロントエンドのサンプルとして `./app-ui/` に [Vite](https://vite.dev) で作成したテンプレートを置いています。

他のframeworkを使用したい場合、フロントエンドのプロジェクトを削除し作成しなおしてください。

`static.go` はGoのバックエンドにフロントエンドの成果物を埋め込むために使用しています。
フロントエンドプロジェクトの名前に応じて適宜変更してください。

```sh
cd ./frontend

rm -r ./app-ui

# See https://vite.dev/guide/
npm create vite@latest
```
