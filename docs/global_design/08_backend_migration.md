# バックエンド移行設計書（Next.js → Go API）

## 📋 概要

このドキュメントでは、現在Next.jsから直接PostgreSQLに接続している実装を、**Go APIサーバー経由でのデータアクセスに移行する**際の設計方針を定義します。

## 🎯 移行の目的

- **責務の分離**: フロントエンド（Next.js）とバックエンド（Go）を明確に分離
- **スケーラビリティ**: Go APIサーバーを独立してスケールさせる
- **パフォーマンス**: Goの高速な処理性能を活用
- **保守性**: フロントエンドとバックエンドを独立して開発・デプロイ可能に

## 🏗️ アーキテクチャ変更

### 現在のアーキテクチャ（Phase 1）

```
Next.js (Frontend)
├─ app/              # App Router
├─ features/         # 機能モジュール
├─ shared/           # 共通コンポーネント
└─ external/         # 外部連携層
   ├─ dto/           # API契約（型定義）
   ├─ handler/       # エントリーポイント
   ├─ service/       # ビジネスロジック
   ├─ domain/        # ドメインインターフェース
   ├─ repository/    # DB直接アクセス ★削除予定
   └─ client/        # DB接続（Drizzle） ★削除予定
       ↓
   PostgreSQL (Neon)
```

### 移行後のアーキテクチャ（Phase 2）

```
Next.js (Frontend)
├─ app/              # App Router
├─ features/         # 機能モジュール
├─ shared/           # 共通コンポーネント
└─ external/         # 外部連携層
   ├─ dto/           # API契約（型定義）★維持
   ├─ handler/       # エントリーポイント ★維持
   ├─ service/       # GoのAPI呼び出し ★変更
   └─ client/        # HTTPクライアント（axios） ★変更
       ↓ HTTP/REST
   Go API Server
   ├─ handler/       # HTTPハンドラー
   ├─ service/       # ビジネスロジック
   ├─ domain/        # ドメインモデル
   ├─ repository/    # DBアクセス
   └─ client/        # DB接続
       ↓
   PostgreSQL (Neon)
```

## 📂 移行対象の詳細

### 削除される層（Next.js側）

以下のディレクトリ・ファイルは**Go APIに移行**されます：

```
external/
├─ domain/           ❌ 削除（Go側に移行）
│  ├─ note/
│  ├─ template/
│  └─ account/
├─ repository/       ❌ 削除（Go側に移行）
│  ├─ note.repository.ts
│  ├─ template.repository.ts
│  └─ account.repository.ts
├─ client/           ❌ 削除（Go側に移行）
│  └─ db.ts          # Drizzle ORM接続
└─ dto/              ❌ 削除（TypeSpecから自動生成に移行）★重要
   ├─ note.dto.ts
   ├─ template.dto.ts
   └─ account.dto.ts
```

**理由**:
- ドメインロジックとDBアクセスはバックエンド（Go）の責務
- Next.jsはフロントエンドとしての役割に専念
- **DTOはTypeSpecから自動生成するため、手書きのDTOは削除**

### 維持される層（Next.js側）

```
external/
├─ generated/        ✨ 新規（TypeSpecから自動生成）
│  ├─ models/        # リクエスト/レスポンス型
│  └─ client/        # API Clientコード
├─ handler/          ✅ 維持（エントリーポイント）
│  ├─ note.query.server.ts
│  ├─ note.command.action.ts
│  └─ ...
├─ service/          ✅ 維持（Go API呼び出しに変更）
│  ├─ note.service.ts
│  ├─ template.service.ts
│  └─ account.service.ts
└─ client/           ✅ 変更（HTTPクライアントに変更）
   └─ api-client.ts  # axios/fetch（手書き or 生成）
```

## 🔄 移行の段階的アプローチ

### Phase 1: 現在（Next.js + DB直接接続）

```ts
// external/service/note.service.ts
export class NoteService {
  constructor(
    private noteRepository: INoteRepository,  // ← Drizzle ORM
    private transactionManager: ITransactionManager<DbClient>,
  ) {}

  async createNote(ownerId: string, input: CreateNoteRequest): Promise<Note> {
    return this.transactionManager.execute(async (tx) => {
      const template = await this.templateRepository.findById(input.templateId, tx);
      return this.noteRepository.create({ title: input.title, ... }, tx);
    });
  }
}
```

### Phase 2: 移行後（Next.js + Go API）

```ts
// external/service/note.service.ts
export class NoteService {
  constructor(
    private apiClient: ApiClient,  // ← HTTPクライアント
  ) {}

  async createNote(ownerId: string, input: CreateNoteRequest): Promise<Note> {
    // Go APIのエンドポイントを呼び出す
    const response = await this.apiClient.post<NoteResponse>('/api/notes', {
      ...input,
      ownerId,
    });

    // DTOで型を保証して返す
    return NoteResponseSchema.parse(response.data);
  }
}
```

**変更点**:
- Repositoryへの直接アクセスを削除
- HTTPクライアント経由でGo APIを呼び出す
- DTOは変更なし（API契約を維持）

## 🔌 Go API側の実装

### エンドポイント設計

Go APIは、`docs/global_design/07_api_design.md`で定義されたAPI仕様を実装します。

#### Notes API

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/notes` | ノート一覧取得 |
| GET | `/api/notes/:id` | ノート詳細取得 |
| POST | `/api/notes` | ノート作成 |
| PUT | `/api/notes/:id` | ノート更新 |
| POST | `/api/notes/:id/publish` | ノート公開 |
| POST | `/api/notes/:id/unpublish` | ノート公開取り消し |
| DELETE | `/api/notes/:id` | ノート削除 |

#### Templates API

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/templates` | テンプレート一覧取得 |
| GET | `/api/templates/:id` | テンプレート詳細取得 |
| POST | `/api/templates` | テンプレート作成 |
| PUT | `/api/templates/:id` | テンプレート更新 |
| DELETE | `/api/templates/:id` | テンプレート削除 |

#### Accounts API

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/accounts/me` | 現在のユーザー取得 |
| GET | `/api/accounts/:id` | アカウント詳細取得 |
| POST | `/api/accounts/auth` | OAuth連携 |

### Go API ディレクトリ構成（参考）

```
backend/
├─ cmd/
│  └─ api/
│     └─ main.go          # エントリーポイント
├─ internal/
│  ├─ handler/            # HTTPハンドラー
│  │  ├─ note_handler.go
│  │  ├─ template_handler.go
│  │  └─ account_handler.go
│  ├─ service/            # ビジネスロジック
│  │  ├─ note_service.go
│  │  ├─ template_service.go
│  │  └─ account_service.go
│  ├─ domain/             # ドメインモデル
│  │  ├─ note/
│  │  │  ├─ note.go
│  │  │  ├─ section.go
│  │  │  └─ repository.go  # interface
│  │  ├─ template/
│  │  └─ account/
│  ├─ repository/         # DBアクセス実装
│  │  ├─ note_repository.go
│  │  ├─ template_repository.go
│  │  └─ account_repository.go
│  └─ client/             # DB接続
│     └─ postgres.go
├─ pkg/
│  ├─ dto/                # リクエスト/レスポンス型
│  └─ errors/             # エラー定義
└─ migrations/            # DBマイグレーション
```

## 🔐 認証の扱い

### 現在（NextAuth）

- Next.jsがNextAuthでGoogle OAuthを処理
- セッション情報をNext.jsが管理

### 移行後

**Option 1: Next.jsで認証、JWTをGo APIに渡す**

```
1. ユーザーがNext.jsにアクセス
2. NextAuthでGoogle認証
3. Next.js → Go API呼び出し時にJWTを付与
4. Go APIがJWTを検証して処理
```

```ts
// external/client/api-client.ts
export const apiClient = axios.create({
  baseURL: process.env.API_BASE_URL,
});

apiClient.interceptors.request.use(async (config) => {
  const session = await getServerSession(authOptions);
  if (session?.accessToken) {
    config.headers.Authorization = `Bearer ${session.accessToken}`;
  }
  return config;
});
```

**Option 2: Go APIで認証を完全に実装**

- Go APIがOAuth処理を担当
- Next.jsはトークンを受け取ってセッション管理のみ

**推奨**: Option 1（Next.jsで認証、JWTでGo APIと連携）

## 🛠️ HTTPクライアントの実装

### Next.js側のAPIクライアント

```ts
// external/client/api-client.ts
import axios, { type AxiosInstance } from 'axios';
import { getServerSession } from 'next-auth';
import { authOptions } from '@/shared/lib/auth';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: process.env.API_BASE_URL || 'http://localhost:8080',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // リクエストインターセプター（認証トークン付与）
    this.client.interceptors.request.use(async (config) => {
      const session = await getServerSession(authOptions);
      if (session?.accessToken) {
        config.headers.Authorization = `Bearer ${session.accessToken}`;
      }
      return config;
    });

    // レスポンスインターセプター（エラーハンドリング）
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // 認証エラー処理
        }
        return Promise.reject(error);
      }
    );
  }

  async get<T>(url: string, params?: Record<string, unknown>) {
    const response = await this.client.get<T>(url, { params });
    return response.data;
  }

  async post<T>(url: string, data?: unknown) {
    const response = await this.client.post<T>(url, data);
    return response.data;
  }

  async put<T>(url: string, data?: unknown) {
    const response = await this.client.put<T>(url, data);
    return response.data;
  }

  async delete<T>(url: string) {
    const response = await this.client.delete<T>(url);
    return response.data;
  }
}

export const apiClient = new ApiClient();
```

## 📊 データフロー比較

### 現在のフロー

```
Client Component
    ↓
Server Action (handler/*.action.ts)
    ↓
Service (external/service/*.service.ts)
    ↓
Repository (external/repository/*.repository.ts)
    ↓
Drizzle ORM (external/client/db.ts)
    ↓
PostgreSQL
```

### 移行後のフロー

```
Client Component
    ↓
Server Action (handler/*.action.ts)
    ↓
Service (external/service/*.service.ts)
    ↓
HTTP Client (external/client/api-client.ts)
    ↓ HTTP/REST
Go API Server
    ↓
Service (internal/service/*.go)
    ↓
Repository (internal/repository/*.go)
    ↓
PostgreSQL Driver
    ↓
PostgreSQL
```

## 🚀 移行手順

### Step 1: Go APIサーバーのセットアップ

1. Goプロジェクトの初期化
2. PostgreSQL接続の実装
3. ドメインモデルの実装（Next.jsのdomainから移植）
4. Repositoryの実装（Next.jsのrepositoryから移植）

### Step 2: API実装（機能ごとに段階的に）

#### 優先順位

1. **Account API** - 認証の基盤
2. **Template API** - 依存が少ない
3. **Note API** - Templateに依存

### Step 3: Next.js側の修正

1. `external/service`の各メソッドをAPI呼び出しに変更
2. `external/client/api-client.ts`の実装
3. 環境変数`API_BASE_URL`の設定

### Step 4: テストと検証

1. Go APIの単体テスト
2. Next.jsからの疎通テスト
3. E2Eテスト

### Step 5: クリーンアップ

1. `external/domain`の削除
2. `external/repository`の削除
3. `external/client/db.ts`の削除
4. Drizzle ORM関連の依存削除

## 🔍 移行チェックリスト

### TypeSpec & コード生成

- [ ] api-schemaリポジトリのセットアップ
- [ ] TypeSpecでモデル定義（Note, Template, Account）
- [ ] TypeSpecでルート定義（エンドポイント）
- [ ] OpenAPI YAML自動生成スクリプト
- [ ] Go型定義自動生成スクリプト
- [ ] TypeScript型定義自動生成スクリプト
- [ ] CI/CDパイプラインの設定（GitHub Actions）

### Go API実装

- [ ] Go APIサーバーのセットアップ
- [ ] PostgreSQL接続の実装
- [ ] 生成された型を使ってAccount APIの実装
- [ ] 生成された型を使ってTemplate APIの実装
- [ ] 生成された型を使ってNote APIの実装
- [ ] トランザクション管理の実装
- [ ] エラーハンドリングの統一

### Next.js側の修正

- [ ] 生成されたTypeScript型のインポート設定
- [ ] Next.js側のHTTPクライアント実装
- [ ] Serviceクラスの書き換え（生成型を使用）
- [ ] 認証フローの実装（JWT連携）
- [ ] 環境変数の設定

### テスト & クリーンアップ

- [ ] Go APIの単体テスト
- [ ] Next.jsからの疎通テスト
- [ ] E2Eテスト
- [ ] `external/domain`の削除
- [ ] `external/repository`の削除
- [ ] `external/client/db.ts`の削除
- [ ] `external/dto`の削除（生成コードに置き換え）
- [ ] Drizzle ORM関連の依存削除
- [ ] ドキュメントの更新

## 🔧 TypeSpec による型定義の自動生成

### 概要

手書きのDTOを廃止し、**TypeSpec**を使ってAPI仕様を定義します。TypeSpecからOpenAPI YAMLを生成し、それを元にGoとTypeScriptの型定義を自動生成します。

### TypeSpecとは

TypeSpecは、Microsoftが開発したAPI定義言語です。OpenAPIやJSON Schemaよりも簡潔で、型安全な定義が可能です。

- 公式サイト: https://typespec.io/
- 特徴: 型安全、簡潔な構文、OpenAPI/JSON Schema自動生成

### コード生成フロー

```
TypeSpec (.tsp)
    ↓ tsp compile
OpenAPI YAML
    ↓ openapi-generator
    ├─ Go (models, client)
    └─ TypeScript (models, client)
```

### ディレクトリ構成の推奨案

#### Option 1: 独立したAPIスキーマリポジトリ（推奨）

```
api-schema/                    # 独立したリポジトリ
├── typespec/
│   ├── main.tsp              # エントリーポイント
│   ├── package.json          # TypeSpec dependencies
│   ├── tspconfig.yaml        # TypeSpec設定
│   ├── models/               # データモデル定義
│   │   ├── note.tsp
│   │   ├── template.tsp
│   │   ├── account.tsp
│   │   └── common.tsp
│   └── routes/               # エンドポイント定義
│       ├── notes.tsp
│       ├── templates.tsp
│       └── accounts.tsp
├── generated/
│   ├── openapi.yaml          # 生成されたOpenAPI仕様
│   ├── go/                   # Go生成コード
│   │   ├── models/
│   │   │   ├── note.go
│   │   │   ├── template.go
│   │   │   └── account.go
│   │   └── client/
│   │       └── api_client.go
│   └── typescript/           # TypeScript生成コード
│       ├── models/
│       │   ├── note.ts
│       │   ├── template.ts
│       │   └── account.ts
│       └── client/
│           └── api-client.ts
├── scripts/
│   ├── generate.sh           # 全自動生成スクリプト
│   ├── generate-openapi.sh   # OpenAPI生成
│   ├── generate-go.sh        # Go生成
│   └── generate-ts.sh        # TypeScript生成
├── package.json
└── README.md
```

**メリット**:
- フロントエンドとバックエンドから独立
- API仕様が単一の真実の源（Single Source of Truth）
- バージョン管理が容易
- 生成コードをnpmパッケージやGoモジュールとして配布可能

#### Option 2: モノレポ内に配置

```
root/
├── frontend/
├── backend/
└── packages/
    └── api-schema/
        ├── typespec/
        ├── generated/
        └── package.json
```

**メリット**:
- 1つのリポジトリで管理
- 依存関係の同期が簡単
- CI/CDパイプラインの統一

**推奨**: **Option 1（独立したリポジトリ）**

理由:
- API仕様はフロントエンドとバックエンドの契約なので、両方から独立すべき
- バージョン管理が明確（セマンティックバージョニング）
- 複数のクライアント（Web、モバイル）で共有しやすい

### TypeSpec定義例

#### モデル定義

```typespec
// typespec/models/note.tsp
import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

namespace MiniNotion.Models;

/** ノートのステータス */
enum NoteStatus {
  Draft: "Draft",
  Publish: "Publish",
}

/** セクション（ノートの各項目） */
model Section {
  /** セクションID */
  id: string;

  /** フィールドID */
  fieldId: string;

  /** フィールドラベル */
  fieldLabel: string;

  /** 内容 */
  content: string;

  /** 必須項目かどうか */
  isRequired: boolean;
}

/** ノート作成リクエスト */
model CreateNoteRequest {
  /** タイトル */
  @minLength(1)
  @maxLength(100)
  title: string;

  /** テンプレートID */
  @format("uuid")
  templateId: string;

  /** セクション（オプション） */
  sections?: Section[];
}

/** ノートレスポンス */
model NoteResponse {
  /** ノートID */
  id: string;

  /** タイトル */
  title: string;

  /** テンプレートID */
  templateId: string;

  /** テンプレート名 */
  templateName: string;

  /** 所有者ID */
  ownerId: string;

  /** 所有者情報 */
  owner: {
    id: string;
    firstName: string;
    lastName: string;
    thumbnail?: string;
  };

  /** ステータス */
  status: NoteStatus;

  /** セクション */
  sections: Section[];

  /** 作成日時 */
  createdAt: utcDateTime;

  /** 更新日時 */
  updatedAt: utcDateTime;
}
```

#### エンドポイント定義

```typespec
// typespec/routes/notes.tsp
import "@typespec/http";
import "@typespec/openapi3";
import "../models/note.tsp";

using TypeSpec.Http;
using MiniNotion.Models;

namespace MiniNotion.Routes;

@route("/api/notes")
interface Notes {
  /** ノート一覧取得 */
  @get
  @summary("Get notes list")
  list(
    /** タイトルキーワード検索 */
    @query q?: string,

    /** ステータスフィルター */
    @query status?: NoteStatus,

    /** テンプレートIDフィルター */
    @query templateId?: string,

    /** 所有者IDフィルター */
    @query ownerId?: string,
  ): NoteResponse[];

  /** ノート詳細取得 */
  @get
  @route("/{noteId}")
  @summary("Get note by ID")
  get(
    @path noteId: string
  ): NoteResponse | NotFoundError;

  /** ノート作成 */
  @post
  @summary("Create note")
  create(
    @body request: CreateNoteRequest
  ): NoteResponse;

  /** ノート更新 */
  @put
  @route("/{noteId}")
  @summary("Update note")
  update(
    @path noteId: string,
    @body request: UpdateNoteRequest
  ): NoteResponse;

  /** ノート公開 */
  @post
  @route("/{noteId}/publish")
  @summary("Publish note")
  publish(
    @path noteId: string
  ): NoteResponse;

  /** ノート公開取り消し */
  @post
  @route("/{noteId}/unpublish")
  @summary("Unpublish note")
  unpublish(
    @path noteId: string
  ): NoteResponse;

  /** ノート削除 */
  @delete
  @route("/{noteId}")
  @summary("Delete note")
  delete(
    @path noteId: string
  ): { success: boolean };
}

/** Not Found エラー */
@error
model NotFoundError {
  code: "NOT_FOUND";
  message: string;
}
```

### 自動生成スクリプト

#### package.json

```json
{
  "name": "@mini-notion/api-schema",
  "version": "1.0.0",
  "scripts": {
    "generate": "npm run generate:openapi && npm run generate:go && npm run generate:ts",
    "generate:openapi": "tsp compile typespec --emit @typespec/openapi3",
    "generate:go": "./scripts/generate-go.sh",
    "generate:ts": "./scripts/generate-ts.sh"
  },
  "devDependencies": {
    "@typespec/compiler": "^0.60.0",
    "@typespec/http": "^0.60.0",
    "@typespec/openapi3": "^0.60.0",
    "@openapitools/openapi-generator-cli": "^2.13.0"
  }
}
```

#### TypeSpec設定（tspconfig.yaml）

```yaml
emit:
  - "@typespec/openapi3"

options:
  "@typespec/openapi3":
    output-file: "generated/openapi.yaml"
```

#### Go生成スクリプト（scripts/generate-go.sh）

```bash
#!/bin/bash

openapi-generator-cli generate \
  -i generated/openapi.yaml \
  -g go \
  -o generated/go \
  --package-name api \
  --git-repo-id mini-notion-api \
  --git-user-id your-org \
  --additional-properties=enumClassPrefix=true
```

#### TypeScript生成スクリプト（scripts/generate-ts.sh）

```bash
#!/bin/bash

openapi-generator-cli generate \
  -i generated/openapi.yaml \
  -g typescript-axios \
  -o generated/typescript \
  --additional-properties=withSeparateModelsAndApi=true,apiPackage=client,modelPackage=models
```

### 生成されるコード例

#### Go側

```go
// generated/go/models/note.go
package models

import "time"

type NoteStatus string

const (
    NoteStatusDraft   NoteStatus = "Draft"
    NoteStatusPublish NoteStatus = "Publish"
)

type CreateNoteRequest struct {
    Title      string     `json:"title" validate:"required,min=1,max=100"`
    TemplateID string     `json:"templateId" validate:"required,uuid"`
    Sections   []Section  `json:"sections,omitempty"`
}

type NoteResponse struct {
    ID           string      `json:"id"`
    Title        string      `json:"title"`
    TemplateID   string      `json:"templateId"`
    TemplateName string      `json:"templateName"`
    OwnerID      string      `json:"ownerId"`
    Owner        Owner       `json:"owner"`
    Status       NoteStatus  `json:"status"`
    Sections     []Section   `json:"sections"`
    CreatedAt    time.Time   `json:"createdAt"`
    UpdatedAt    time.Time   `json:"updatedAt"`
}
```

#### TypeScript側

```typescript
// generated/typescript/models/note.ts
export enum NoteStatus {
    Draft = 'Draft',
    Publish = 'Publish'
}

export interface Section {
    id: string;
    fieldId: string;
    fieldLabel: string;
    content: string;
    isRequired: boolean;
}

export interface CreateNoteRequest {
    title: string;
    templateId: string;
    sections?: Section[];
}

export interface NoteResponse {
    id: string;
    title: string;
    templateId: string;
    templateName: string;
    ownerId: string;
    owner: {
        id: string;
        firstName: string;
        lastName: string;
        thumbnail?: string;
    };
    status: NoteStatus;
    sections: Section[];
    createdAt: string;
    updatedAt: string;
}
```

### Next.js側での使用方法

生成されたTypeScript型を使用します。

```ts
// external/service/note.service.ts
import { NoteResponse, CreateNoteRequest } from '@mini-notion/api-schema/typescript/models';
import { apiClient } from '../client/api-client';

export class NoteService {
  async createNote(ownerId: string, input: CreateNoteRequest): Promise<NoteResponse> {
    // 生成された型を使用
    const response = await apiClient.post<NoteResponse>('/api/notes', {
      ...input,
      ownerId,
    });

    return response;
  }

  async listNotes(filters?: {
    q?: string;
    status?: NoteStatus;
    templateId?: string;
  }): Promise<NoteResponse[]> {
    return apiClient.get<NoteResponse[]>('/api/notes', filters);
  }
}
```

### Go側での使用方法

```go
// internal/handler/note_handler.go
import (
    "github.com/your-org/mini-notion-api/generated/go/models"
)

func (h *NoteHandler) CreateNote(c *gin.Context) {
    var req models.CreateNoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    note, err := h.service.CreateNote(c, userID, &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, note)
}
```

### CI/CDパイプライン

#### GitHub Actions例

```yaml
# api-schema/.github/workflows/generate.yml
name: Generate API Code

on:
  push:
    branches: [main]
    paths:
      - 'typespec/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: npm install

      - name: Generate OpenAPI
        run: npm run generate:openapi

      - name: Generate Go code
        run: npm run generate:go

      - name: Generate TypeScript code
        run: npm run generate:ts

      - name: Commit generated files
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add generated/
          git commit -m "chore: regenerate API code" || exit 0
          git push
```

### TypeSpec移行のメリット

1. **型の一貫性**: GoとTypeScriptで同じ型定義を使用
2. **変更の追跡**: TypeSpecを変更すれば、両方のコードが自動更新
3. **ドキュメント生成**: OpenAPIから自動でドキュメント生成（Swagger UI等）
4. **バリデーション**: TypeSpecレベルで型チェック
5. **保守性**: 手書きDTOのメンテナンス不要

### トランザクション管理

現在Next.jsで実装しているトランザクション管理は、Go API側に移行します。

```go
// internal/service/note_service.go
func (s *NoteService) CreateNote(ctx context.Context, ownerID string, req *dto.CreateNoteRequest) (*dto.NoteResponse, error) {
    return s.txManager.Execute(ctx, func(tx *sql.Tx) (*dto.NoteResponse, error) {
        // トランザクション内での処理
        template, err := s.templateRepo.FindByID(ctx, req.TemplateID, tx)
        if err != nil {
            return nil, err
        }

        note, err := s.noteRepo.Create(ctx, &domain.Note{
            Title:      req.Title,
            TemplateID: req.TemplateID,
            OwnerID:    ownerID,
        }, tx)

        return toNoteResponse(note), nil
    })
}
```

## 🌐 環境変数

### Next.js（`.env.local`）

```bash
# Go APIのベースURL
API_BASE_URL=http://localhost:8080

# 本番環境
# API_BASE_URL=https://api.mini-notion.com
```

### Go API（`.env`）

```bash
# PostgreSQL接続
DATABASE_URL=postgresql://user:password@localhost:5432/mini_notion

# サーバー設定
PORT=8080
ENV=development

# 認証
JWT_SECRET=your-secret-key
```

## 📚 参考資料

- `docs/global_design/07_api_design.md` - API仕様
- `frontend/docs/05_external_layer.md` - External Layerの設計
- `docs/global_design/05_domain_design.md` - ドメイン設計
- `docs/global_design/06_database_design.md` - データベース設計
