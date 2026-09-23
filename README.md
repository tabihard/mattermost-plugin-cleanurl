# mattermost-plugin-template

Mattermost の公式チュートリアル [Hello World (server plugin)](https://developers.mattermost.com/integrate/plugins/components/server/hello-world/) を元にした、サーバープラグイン用のテンプレートリポジトリです。

`/plugins/com.tabihard.mattermost-plugin-template` にアクセスすると `Hello, world!` を返すだけの最小構成のプラグインを、Go modules と `Makefile` でビルド・パッケージングできる状態にしてあります。

## 構成

```
.
├── go.mod
├── plugin.json          # プラグインのマニフェスト
├── server/
│   └── main.go           # プラグイン本体 (ServeHTTP で "Hello, world!" を返す)
└── Makefile              # クロスコンパイル & tar.gz バンドル生成
```

## 前提条件

- Go 1.23 以上（`go.mod` の `toolchain` 指定により、依存先 `mattermost/server/public` が要求する Go 1.26.7 が `GOTOOLCHAIN=auto` で自動取得されます）
- プラグインを有効化した Mattermost サーバー（`PluginSettings > Enable = true`）

## ビルド

```bash
go mod tidy   # 依存関係の取得・go.sum の生成（初回のみ）
make          # dist/com.tabihard.mattermost-plugin-template.tar.gz を生成
```

`make build` だけを実行すると、`linux-amd64` / `linux-arm64` / `darwin-amd64` / `darwin-arm64` / `windows-amd64` 向けのバイナリを `dist/<plugin id>/server/dist/` 配下に生成します（`plugin.json` の `server.executables` に対応）。

## インストール

**System Console から:**

1. System Admin でログインし `/admin_console` を開く
2. *Plugins > Plugin Management* で `dist/com.tabihard.mattermost-plugin-template.tar.gz` をアップロード
3. アップロード後にプラグインを Enable にする

**手動インストール:**

`dist/com.tabihard.mattermost-plugin-template/` の中身を、Mattermost サーバーの `PluginSettings > Directory`（デフォルトは `./plugins`）配下の `plugins/com.tabihard.mattermost-plugin-template/` に配置し、サーバーを再起動します。

## 動作確認

```
https://<your-mattermost-server>/plugins/com.tabihard.mattermost-plugin-template
```

にアクセスして `Hello, world!` が表示されれば成功です。

## このテンプレートの使い方

1. リポジトリを clone/fork する
2. `go.mod` の module path、`plugin.json` の `id` / `name` / `homepage_url` などを自分のプラグイン用に書き換える
3. `server/main.go` に実装を追加していく

より本格的な構成（Webapp 側の React コンポーネント、CI、lint 設定など）が必要な場合は、公式の [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template) も参考にしてください。
