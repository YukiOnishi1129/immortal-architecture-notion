# API設計書（MVP）

## 概要

このドキュメントは、システムが提供するAPIの仕様を定義します。実装詳細については、フロントエンド層のドキュメント（`frontend/docs/`）を参照してください。

## API構成

### エンドポイント分類

- **Query（読み取り）**: データ取得のみ。副作用なし（GET）
- **Command（書き込み）**: データの作成・更新・削除。副作用あり（POST, PUT, DELETE）

### URL設計とHTTPメソッド

| 操作 | HTTPメソッド | URLパターン | 用途 |
|------|-------------|------------|------|
| 一覧取得 | GET | `/api/xxx` | 全件または条件付き取得（クエリパラメータで絞り込み） |
| 単体取得 | GET | `/api/xxx/:id` | IDで1件取得 |
| 作成 | POST | `/api/xxx` | 新規作成 |
| 更新 | PUT | `/api/xxx/:id` | 既存更新 |
| 削除 | DELETE | `/api/xxx/:id` | 削除 |
| 状態変更 | POST | `/api/xxx/:id/action` | 状態遷移（例: `/api/notes/:id/publish`） |

---

## Notes（ノート）API

### Query Operations

#### ノート一覧取得

**URL**: `GET /api/notes`

**Request (Query Parameters)**:
```
NoteFilters {
  q?: string                    // タイトルのキーワード検索
  status?: "Draft" | "Publish"  // ステータスフィルター
  templateId?: string           // テンプレートIDフィルター
  ownerId?: string              // 所有者IDでフィルタ（自分のノートのみ取得する場合に使用）
}
```

**Response**:
```
NoteResponse {
  id: string
  title: string
  templateId: string
  templateName: string
  ownerId: string
  owner: {
    id: string
    firstName: string
    lastName: string
    thumbnail: string?
  }
  status: "Draft" | "Publish"
  sections: [{
    id: string
    fieldId: string
    fieldLabel: string
    content: string
    isRequired: boolean
  }]
  createdAt: string  // ISO 8601形式
  updatedAt: string  // ISO 8601形式
}

ListNoteResponse = NoteResponse[]
```

**ビジネスルール**:
- 認証必須
- 公開済み（Publish）のノートまたは自分のノートを取得可能
- `ownerId`を指定した場合、そのユーザーが所有するノートのみを取得
- 自分のノートのみを取得する場合: `GET /api/notes?ownerId={自分のID}`

---

#### ノート詳細取得

**URL**: `GET /api/notes/:id`

**Request (URL Parameters)**:
```
id: string  // ノートID
```

**Response**:
```
GetNoteByIdResponse = NoteResponse | null;  // 見つからない場合はnull
```

**ビジネスルール**:
- 認証必須
- 存在しないIDの場合はnullを返す

---

### Command Operations

#### ノート作成

**URL**: `POST /api/notes`

**Request**:
```
CreateNoteRequest {
  title: string
  templateId: string
  sections: [{
    fieldId: string
    content: string
  }]
}
```

**Response**:
```
CreateNoteResponse = NoteResponse;
```

**ビジネスルール**:
- 認証必須
- 新規作成時のステータスは"Draft"
- 指定されたテンプレートが存在する必要がある
- sectionsは必須（テンプレートの全フィールドに対応するセクションが必要）
- isRequiredがtrueのフィールドはcontentが空だとエラー

---

#### ノート更新

**URL**: `PUT /api/notes/:id`

**Request**:
```
UpdateNoteRequest {
  id: string       // ノートID
  title: string
  sections: [{
    id: string     // セクションID
    content: string
  }>;
}
```

**Response**:
```
UpdateNoteResponse = NoteResponse;
```

**ビジネスルール**:
- 認証必須
- 自分が所有するノートのみ更新可能
- テンプレートのフィールド構造は変更不可

---

#### ノート公開

**URL**: `POST /api/notes/:id/publish`

**Request**:
```
PublishNoteRequest {
  noteId: string
}
```

**Response**:
```
PublishNoteResponse = NoteResponse;
```

**ビジネスルール**:
- 認証必須
- 自分が所有するノートのみ公開可能
- 下書き（Draft）から公開済み（Publish）に状態遷移
- 既に公開済みの場合はエラー

---

#### ノート公開取り消し

**URL**: `POST /api/notes/:id/unpublish`

**Request**:
```
UnpublishNoteRequest {
  noteId: string
}
```

**Response**:
```
UnpublishNoteResponse = NoteResponse;
```

**ビジネスルール**:
- 認証必須
- 自分が所有するノートのみ公開取り消し可能
- 公開済み（Publish）から下書き（Draft）に状態遷移
- 既に下書きの場合はエラー

---

#### ノート削除

**URL**: `DELETE /api/notes/:id`

**Request**:
```
DeleteNoteRequest {
  id: string  // ノートID
}
```

**Response**:
```
DeleteNoteResponse {
  success: boolean
}
```

**ビジネスルール**:
- 認証必須
- 自分が所有するノートのみ削除可能
- ノートに紐づくセクションも同時に削除される

---

## Templates（テンプレート）API

### Query Operations

#### テンプレート一覧取得

**URL**: `GET /api/templates`

**Request (Query Parameters)**:
```
q?: string         // テンプレート名のキーワード検索
ownerId?: string   // 所有者IDでフィルタ（自分のテンプレートのみ取得する場合に使用）
```

**Response**:
```
TemplateResponse {
  id: string
  name: string
  ownerId: string
  owner: {
    id: string
    firstName: string
    lastName: string
    thumbnail: string?;
  };
  fields: [{
    id: string
    label: string
    order: number
    isRequired: boolean
  }>;
  updatedAt: string  // ISO 8601形式
  isUsed: boolean    // ノートで使用中かどうか
}

ListTemplatesResponse = TemplateResponse[];
```

**ビジネスルール**:
- 認証必須
- `ownerId`を指定した場合、そのユーザーが所有するテンプレートのみを取得
- 自分のテンプレートのみを取得する場合: `GET /api/templates?ownerId={自分のID}`
- `isUsed`は、テンプレートがノートで使用中かを示す

---

#### テンプレート詳細取得

**URL**: `GET /api/templates/:id`

**Request (URL Parameters)**:
```
id: string  // テンプレートID
```

**Response**:
```
GetTemplateByIdResponse = TemplateResponse | null;  // 見つからない場合はnull
```

**ビジネスルール**:
- 認証必須
- 存在しないIDの場合はnullを返す

---

### Command Operations

#### テンプレート作成

**URL**: `POST /api/templates`

**Request**:
```
CreateTemplateRequest {
  name: string
  fields: [{
    label: string
    order: number
    isRequired: boolean
  }>;
}
```

**Response**:
```
CreateTemplateResponse = TemplateResponse;
```

**ビジネスルール**:
- 認証必須
- フィールドのorderは0から始まる連番
- 新規作成時のisUsedはfalse

---

#### テンプレート更新

**URL**: `PUT /api/templates/:id`

**Request**:
```
UpdateTemplateRequest {
  id: string       // テンプレートID
  name: string
  fields: [{
    id?: string    // 既存フィールドの場合は必須
    label: string
    order: number
    isRequired: boolean
  }>;
}
```

**Response**:
```
UpdateTemplateResponse = TemplateResponse;
```

**ビジネスルール**:
- 認証必須
- 自分が所有するテンプレートのみ更新可能
- **テンプレートがノートで使用中（isUsed = true）の場合**:
  - テンプレート名の変更: 可能
  - フィールドのlabel変更: 可能
  - フィールドのisRequired変更: 可能
  - フィールドの追加: 不可
  - フィールドの削除: 不可
  - フィールドのorder変更: 不可
- **テンプレートが未使用（isUsed = false）の場合**:
  - すべての変更が可能

---

#### テンプレート削除

**URL**: `DELETE /api/templates/:id`

**Request**:
```
DeleteTemplateRequest {
  id: string  // テンプレートID
}
```

**Response**:
```
DeleteTemplateResponse {
  success: boolean
}
```

**ビジネスルール**:
- 認証必須
- 自分が所有するテンプレートのみ削除可能
- ノートで使用中（isUsed = true）のテンプレートは削除不可
- テンプレートに紐づくフィールドも同時に削除される

---

## Accounts（アカウント）API

### OAuth連携時のアカウント作成または取得

**URL**: `POST /api/accounts/auth` (内部処理)

**Request**:
```
CreateOrGetAccountRequest {
  email: string
  name: string
  provider: string           // 例: "google"
  providerAccountId: string
  thumbnail?: string
}
```

**Response**:
```
AccountResponse {
  id: string
  email: string
  firstName: string
  lastName: string
  fullName: string
  thumbnail: string?;
  lastLoginAt: string  // ISO 8601形式
  createdAt: string    // ISO 8601形式
  updatedAt: string    // ISO 8601形式
}
```

**ビジネスルール**:
- 既存アカウントが存在する場合は取得、存在しない場合は新規作成
- nameは姓名に分割される

---

### 現在のアカウント取得

**URL**: `GET /api/accounts/me`

**Request**: なし

**Response**:
```
GetCurrentAccountResponse = AccountResponse;
```

**ビジネスルール**:
- 認証必須
- ログインユーザーのアカウント情報を取得

---

### アカウント詳細取得

**URL**: `GET /api/accounts/:id`

**Request (URL Parameters)**:
```
id: string  // アカウントID
```

**Response**:
```
GetAccountByIdResponse = AccountResponse | null;  // 見つからない場合はnull
```

**ビジネスルール**:
- 認証必須
- 存在しないIDの場合はnullを返す

---

## ドメインモデルの関係

### エンティティの関連

```
Account (アカウント)
  |
  +-- Template (テンプレート)
  |     |
  |     +-- Field (フィールド)
  |
  +-- Note (ノート)
        |
        +-- Section (セクション)
```

### 関係性の説明

- **Account**: システムのユーザーを表す
- **Template**: ノートの構造を定義する
  - 1つのTemplateは複数のFieldを持つ
  - 1つのAccountは複数のTemplateを所有できる
- **Field**: Templateの項目を定義する
  - label（ラベル）、order（順序）、isRequired（必須フラグ）を持つ
- **Note**: ユーザーが作成するコンテンツ
  - 1つのTemplateに基づいて作成される
  - 1つのAccountが所有する
  - 1つのNoteは複数のSectionを持つ
- **Section**: Noteの各項目の内容
  - Templateのfieldに対応する
  - 実際のコンテンツを保持する

---

## 認証・認可の方針

### 認証方式

- **Google OAuth 2.0**による認証
- 認証チェックはフロント（Next.js BE）で実施し、バックエンド API は信頼できるクライアントからのみ呼ばれる前提。バックエンド側ではヘッダーでの認証検証は行わず、リクエストで渡された `ownerId` をそのまま用いる。

### 認可（権限チェック）

#### 1. Ownerチェック

- リソースの所有者のみが操作可能
- 適用対象:
  - ノートの更新・削除・公開・公開取り消し
  - テンプレートの更新・削除

#### 2. ステータスベースの制御

**ノート**:
- 公開（Publish）: すべてのユーザーが閲覧可能
- 下書き（Draft）: 所有者のみが閲覧可能

**テンプレート**:
- 使用中（isUsed = true）: フィールド構造の変更不可
- 未使用（isUsed = false）: すべての変更が可能

### 権限チェックの考え方

| 操作 | 認証 | Owner確認 | その他の条件 |
|-----|------|----------|------------|
| ノート一覧取得 | 必須 | 不要（ownerIdでフィルタ可） | 公開済みまたは自分のノート |
| ノート詳細取得 | 必須 | 不要 | 公開済みまたは自分のノート |
| ノート作成 | 必須 | 自動設定 | - |
| ノート更新 | 必須 | 必須 | - |
| ノート公開 | 必須 | 必須 | Draft状態のみ |
| ノート公開取り消し | 必須 | 必須 | Publish状態のみ |
| ノート削除 | 必須 | 必須 | - |
| テンプレート一覧取得 | 必須 | 不要（ownerIdでフィルタ可） | - |
| テンプレート詳細取得 | 必須 | 不要 | - |
| テンプレート作成 | 必須 | 自動設定 | - |
| テンプレート更新 | 必須 | 必須 | 使用中の場合は制限あり |
| テンプレート削除 | 必須 | 必須 | 未使用のみ |

---

## 型定義の補足

### 共通型

```
// ノートのステータス
NoteStatus = "Draft" | "Publish";

// 日付形式
ISODateString = string;  // ISO 8601形式（例: "2025-11-16T09:00:00Z"）
```

### バリデーションルール（概念）

- **title**: 1文字以上の文字列
- **name**: 1文字以上の文字列
- **label**: 1文字以上の文字列
- **content**: 0文字以上の文字列（空文字可）
- **order**: 0以上の整数
- **isRequired**: boolean
- **id**: UUID v4形式の文字列
