# アーキテクチャ設計

## 全体構成

```
frontend/src/
├─ app/          # App Router (薄く保つ)
├─ features/     # 機能別モジュール
├─ shared/       # 共通コンポーネント・ユーティリティ
└─ external/     # 外部連携層 (API・DB)
```

## 設計原則

1. **関心の分離**: 各層の責任を明確に定義
2. **Server Components優先**: クライアントサイドのJSを最小化
3. **型安全性**: TypeScriptとZodによる完全な型保証
4. **テスタビリティ**: 各層を独立してテスト可能に
5. **変更可用性**: バックエンド技術の変更に対する柔軟性を確保

## レイヤーの責務

### App Router (`/app`)
- ルーティング定義
- メタデータ設定
- 認証チェック
- エラーハンドリング

### Features (`/features`)
- ビジネスロジック
- UI実装
- 状態管理
- カスタムフック

### Shared (`/shared`)
- 共通コンポーネント
- ユーティリティ関数
- 型定義
- プロバイダー

### External (`/external`)
- データアクセス（現在）/ API連携（将来）
- ビジネスロジック実装
- データ変換（DTO）
- 変更可用性の確保

## データフロー

```mermaid
graph TD
    A[Page Component] --> B[Feature Template]
    B --> C[Container Component]
    C --> D[Custom Hook]
    D --> E[Server Action]
    E --> F[Service Layer]
    F --> G[Repository]
    G --> H[Database/API]
```

## 認証アーキテクチャ

- NextAuth.jsによるセッション管理
- ミドルウェアでのルート保護
- Layout componentでの認証状態チェック
- Server-sideでのセッション検証