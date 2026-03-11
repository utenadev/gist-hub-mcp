# Antigravity Mission: PWA Development for gist-hub

## 1. 概要
`gist-hub-mcp` プロジェクトのフロントエンド（PWA）実装を移譲します。
この PWA は、GitHub Gist に保存された「暗号化された開発メモ」を、外出先から安全に閲覧・検索するためのツールです。

## 2. データ契約 (Data Contract)
PWA は、ルートインデックスとなる `index.json` (暗号化済み) を読み込み、各リポジトリのドキュメントを特定します。

### index.json Schema (v1.0)
```json
{
  "managed_repositories": [
    {
      "name": "user/repo-a",
      "index_gist_id": "GIST_ID_A"
    }
  ]
}
```

## 3. 暗号化仕様 (AES-256-GCM)
Go バックエンドと JS フロントエンドで完全な整合性を保つ必要があります。

| パラメータ | 指定値 | 備考 |
| :--- | :--- | :--- |
| アルゴリズム | AES-256-GCM | 改ざん検知 (Auth Tag) を含む |
| 鍵派生 | PBKDF2 | パスフレーズから生成 |
| 反復回数 | 100,000 | SHA-256 使用 |
| Salt | 16 bytes | 保存済みデータから抽出 |
| Nonce (IV) | 12 bytes | 標準 GCM サイズ |

### データ構造 (Binary Layout)
暗号化された Gist コンテンツ（Base64デコード後）の構造は以下の通りです。
`[Salt(16 bytes)] + [Nonce(12 bytes)] + [Ciphertext + AuthTag]`

## 4. テスト用サンプルデータ (Verification Data)
JS 実装の復号テストに使用してください。

- **Passphrase**: `gist-hub-test`
- **Plaintext**: `PWA connection successful!`
- **Encrypted (Base64)**: `ckEf+9u9hjDn13l7cfclDrn6yFvavTyhE/bzK19Ttgqn1fkRBCfEq29Gw5tyuUj5/3vhG6ynlKI31EtOx8S9x8su8F3sLQ==`

## 5. Q&A (Antigravity 2026-03-10)
- **Q: index.json はどこから取得する？**
  - **A**: GitHub Gist API から直接取得してください。PWA はスタンドアロンで動作します。
- **Q: 認証はどうする？**
  - **A**: 初回起動時に GitHub PAT を入力させ、セキュアな場所（セッションストレージ等）に保持してください。
- **Q: データの連結形式は？**
  - **A**: 上記「データ構造」の通りです。

## 6. 期待する機能
1. **Offline-First**: Service Worker と IndexedDB を使用。
2. **Secure Peeking**: パスフレーズを入力するまでデータは一切復号しない。
3. **Search**: 復号されたインデックス内でのインクリメンタル検索。

---
*Organized by Gemini-CLI (Architect) - 2026-03-10*
