# Mini Notion API Schema

TypeSpecを使用したAPI仕様の定義と、TypeScriptの型自動生成を行うプロジェクトです。

## 概要

このリポジトリは、Mini NotionのフロントエンドとバックエンドのAPI契約を定義し、型安全性を保証します。

- **TypeSpec**: API仕様を定義
- **OpenAPI YAML**: TypeSpecから自動生成
- **TypeScript**: OpenAPIからTypeScriptの型とクライアントを自動生成

## ディレクトリ構成

```
api-schema/
├── typespec/              # TypeSpec定義
│   ├── models/           # データモデル定義
│   │   ├── note.tsp
│   │   ├── template.tsp
│   │   ├── account.tsp
│   │   └── common.tsp
│   ├── routes/           # エンドポイント定義
│   │   ├── notes.tsp
│   │   ├── templates.tsp
│   │   └── accounts.tsp
│   ├── main.tsp          # エントリーポイント
│   └── tspconfig.yaml    # TypeSpec設定
├── generated/            # 自動生成アセット
│   └── openapi.yaml      # TypeSpec→OpenAPIの成果物（Go/TSは各プロジェクトへ配置）
├── scripts/              # 生成スクリプト
│   ├── generate.sh
│   ├── generate-openapi.sh
│   └── generate-ts.sh
├── package.json
└── README.md
```

## セットアップ

### 必要な環境

- Node.js 20+
- pnpm 9+

### インストール

```bash
pnpm install
```

## 使い方

### 全自動生成

TypeSpecからOpenAPI YAMLとTypeScriptのコードを一括生成します。

```bash
pnpm run generate
```

### 個別生成

#### OpenAPI YAMLの生成

```bash
pnpm run generate:openapi
```

#### TypeScriptコードの生成

```bash
pnpm run generate:ts
```

## TypeSpec定義の編集

### モデルの追加・編集

`typespec/models/` 配下にTypeSpecファイルを作成または編集します。

例：
```typespec
// typespec/models/note.tsp
model Note {
  id: string;
  title: string;
  status: NoteStatus;
}
```

### ルート（エンドポイント）の追加・編集

`typespec/routes/` 配下にTypeSpecファイルを作成または編集します。

例：
```typespec
// typespec/routes/notes.tsp
@route("/api/notes")
interface Notes {
  @get list(): Note[];
  @post create(@body note: CreateNoteRequest): Note;
}
```

### 生成コードの更新

TypeSpecを編集したら、以下のコマンドで生成コードを更新します。

```bash
pnpm run generate
```

## 生成されたコードの使用

```typescript
import { NoteResponse, CreateNoteRequest } from '@mini-notion/api-schema/typescript/models';

async function createNote(req: CreateNoteRequest): Promise<NoteResponse> {
    // ...
}
```

## CI/CD

GitHub Actionsを使用して、TypeSpecの変更時に自動的にコード生成を実行します。

## ライセンス

MIT
