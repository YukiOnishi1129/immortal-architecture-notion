# 技術スタック

## コア技術

| カテゴリ  | 技術       | バージョン/備考  |
| --------- | ---------- | ---------------- |
| Framework | Next.js    | 16+ (App Router) |
| 言語      | TypeScript | 5.x              |
| UI        | React      | 19               |

## 主要ライブラリ

| 用途              | ライブラリ            | 理由                             |
| ----------------- | --------------------- | -------------------------------- |
| スタイリング      | Tailwind CSS          | ユーティリティファースト CSS     |
| UI コンポーネント | shadcn/ui             | カスタマイズ可能なコンポーネント |
| 状態管理          | TanStack Query        | サーバー状態の管理に特化         |
| フォーム          | React Hook Form + Zod | 型安全なフォーム処理             |
| 認証              | NextAuth              | v4 (Google ログインのみ)         |
| ORM               | Drizzle               | 型安全な SQL                     |
| コード品質        | Biome                 | 高速なリンター/フォーマッター    |

## インフラストラクチャ

| 環境             | サービス         | 詳細                             |
| ---------------- | ---------------- | -------------------------------- |
| 本番 DB          | Neon             | PostgreSQL 互換のサーバーレス DB |
| 開発 DB          | Docker Compose   | PostgreSQL 15                    |
| ホスティング     | Google Cloud Run | コンテナベースのサーバーレス     |
| 認証プロバイダー | Google OAuth     | ソーシャルログインのみ           |

## 開発環境要件

- Node.js 22.x 以上
- pnpm 10.x 以上
- Docker Desktop (開発 DB 用)
