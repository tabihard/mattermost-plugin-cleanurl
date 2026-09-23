# mattermost-plugin-cleanurl

投稿されたメッセージ内のURLから `utm_source` / `fbclid` / `gclid` などのトラッキングパラメータを自動的に削除するMattermostサーバープラグインです。

## 動作

- `MessageWillBePosted` / `MessageWillBeUpdated` フックで投稿・編集時のメッセージ本文をチェックします。
- メッセージ中の `http(s)://` で始まるURLをすべて検出し、既知のトラッキングパラメータだけをクエリ文字列から取り除きます。
- Markdownリンク `[text](https://example.com?utm_source=x)` の括弧などは壊さず、URL部分のみを書き換えます。
- トラッキングパラメータが含まれていないURLや、URLを含まないメッセージは一切変更しません（順序・エンコーディングも保持）。

### デフォルトで削除されるパラメータ

- `utm_` で始まるすべてのパラメータ（`utm_source`, `utm_medium`, `utm_campaign` など）
- `gclid`, `gclsrc`, `dclid`（Google広告系）
- `fbclid`（Facebook）
- `msclkid`（Microsoft広告）
- `mc_cid`, `mc_eid`（Mailchimp）
- `igshid`, `igsh`（Instagram）
- `yclid`（Yandex）
- `twclid`, `ttclid`（Twitter/X, TikTok広告）
- `vero_id`, `mkt_tok`, `_hsenc`, `_hsmi`（マーケティングツール系）
- `ref_src`, `ref_url`, `spm`, `scid`, `si`, `s_kwcid`, `wt.mc_id`

管理コンソールの設定画面（*System Console > Plugins > Clean URL*）から、`ExtraTrackingParams` にカンマ区切りで追加のパラメータ名を指定できます（例: `ref,igshid`、大文字小文字は区別しません）。

## 構成

```
.
├── go.mod
├── plugin.json          # プラグインのマニフェスト・設定スキーマ
├── server/
│   ├── main.go           # フック実装とトラッキングパラメータ除去ロジック
│   └── main_test.go       # ロジックの単体テスト
└── Makefile              # クロスコンパイル & tar.gz バンドル生成
```

## 前提条件

- Go 1.23 以上（`go.mod` の `toolchain` 指定により、依存先 `mattermost/server/public` が要求する Go 1.26.7 が `GOTOOLCHAIN=auto` で自動取得されます）
- プラグインを有効化した Mattermost サーバー（`PluginSettings > Enable = true`）

## ビルド

```bash
go mod tidy   # 依存関係の取得・go.sum の生成（初回のみ）
go test ./... # 単体テストの実行
make          # dist/com.tabihard.mattermost-plugin-cleanurl.tar.gz を生成
```

`make build` だけを実行すると、`linux-amd64` / `linux-arm64` / `darwin-amd64` / `darwin-arm64` / `windows-amd64` 向けのバイナリを `dist/<plugin id>/server/dist/` 配下に生成します（`plugin.json` の `server.executables` に対応）。

## インストール

**System Console から:**

1. System Admin でログインし `/admin_console` を開く
2. *Plugins > Plugin Management* で `dist/com.tabihard.mattermost-plugin-cleanurl.tar.gz` をアップロード
3. アップロード後にプラグインを Enable にする

**手動インストール:**

`dist/com.tabihard.mattermost-plugin-cleanurl/` の中身を、Mattermost サーバーの `PluginSettings > Directory`（デフォルトは `./plugins`）配下の `plugins/com.tabihard.mattermost-plugin-cleanurl/` に配置し、サーバーを再起動します。

## 動作確認

有効化した状態で、チャンネルに `https://example.com/?utm_source=test&id=1` のようなURLを含むメッセージを投稿し、投稿後に `https://example.com/?id=1` へ書き換わっていれば成功です。
