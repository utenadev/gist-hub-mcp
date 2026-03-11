# gist-hub-mcp

[English version (README-en.md)](README-en.md)

`gist-hub-mcp` は、GitHub Gist をバックエンドとする、AI エージェントのための高度な知識管理ハブです。

## 特徴

- **Transparent Encryption**: AES-256-GCM による透過的暗号化をサポート。
- **GistPad 互換**: バックスラッシュ（`\`）セパレータによる階層構造をサポートし、VSCode GistPad と完全互換。
- **SQLite インデックス**: 高速な検索とメタデータ管理のためのローカルキャッシュ。
- **マルチプラットフォーム**: Go CLI (MCP対応) と PWA (Vite + Web Crypto API) の連携。
- **CI/CD パイプライン**: GitHub Actions による自動テストと GoReleaser によるマルチプラットフォーム配布。

## 使い方 (CLI)

GitHub CLI (`gh`) の認証情報を利用します。

```bash
# Gist 一覧を表示 (gist-hub: プレフィックス付きのみ)
gist-hub list

# 暗号化して Gist を作成
gist-hub create <dir> --passphrase "your-password"

# 暗号化された Gist を復号して取得
gist-hub get <id> --passphrase "your-password"
```

## 開発ステータス

- **Phase 1 & 2**: CLI 基礎、GitHub 連携、Cobra リファクタリング (完了 ✅)
- **Phase 3**: SQLite キャッシュ、AES-256-GCM 暗号化統合 (完了 ✅)
- **Phase 4**: Wiki モード（自動インデックス解決）、MCP サーバー化 (進行中 🏗️)

## ライセンス

[MIT License](LICENSE)
