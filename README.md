# gist-hub-mcp

[English version (README-en.md)](README-en.md)

`gist-hub-mcp` は、GitHub Gist をバックエンドとする、AI エージェントのための高度な知識管理ハブです。

## 特徴

- **Universal Library**: GitHub Gist を永続的なナレッジベースとして利用します。
- **Wiki モード**: VSCode GistPad との互換性を持ち、`_index.md` や `_category.json` を用いた階層構造をサポート。
- **SQLite キャッシュ**: API レートリミットを回避し、高速な全文検索を実現。
- **MCP 準拠**: Model Context Protocol (MCP) を通じて、あらゆる AI エージェントに Gist の知識を提供。
- **BBS 連携**: `agent-hub` と連携し、知識の更新をリアルタイムで通知。

## 使い方 (CLI)

現在は開発中の Phase 1 です。GitHub CLI (`gh`) の認証情報を利用します。

```bash
# Gist 一覧を表示 (gist-hub: プレフィックス付きのみ)
gist-hub list

# 特定の Gist を取得
gist-hub get <gist_id>
```

## ライセンス

[MIT License](LICENSE)
