# External Layer (外部連携層)

## 概要

External層は、アプリケーションと外部システム（DB、API）との境界を管理します。**将来的にバックエンドをGoなどの別のAPIに移行する際の変更可用性を考慮**し、Next.jsから直接DBに接続せず、API層を通じてデータをやり取りする設計になっています。

## 設計思想

- **変更可用性**: バックエンドの実装が変わっても、フロントエンドの変更を最小限に抑える
- **関心の分離**: ビジネスロジックと外部依存を明確に分離
- **型安全性**: DTOによる入出力の型保証

## ディレクトリ構造

```
external/
├─ dto/          # データ転送オブジェクト（API契約）
├─ handler/      # エントリーポイント（CQRSパターン）
├─ service/      # ビジネスロジック・API呼び出し
└─ client/       # HTTPクライアント・DB接続（将来的にはAPIクライアントのみ）
```

## レイヤーの責務

### DTO (Data Transfer Object)

APIとの契約を定義。バックエンドがどの技術で実装されていても、このインターフェースは維持されます。

```ts
// external/dto/note.dto.ts
import { z } from 'zod'

// APIリクエストの型定義
export const CreateNoteRequestSchema = z.object({
  title: z.string().min(1).max(100),
  templateId: z.string().uuid(),
})

// APIレスポンスの型定義
export const NoteResponseSchema = z.object({
  id: z.string().uuid(),
  title: z.string(),
  status: z.enum(['Draft', 'Publish']),
  sections: z.array(z.object({
    id: z.string(),
    fieldId: z.string(),
    content: z.string(),
  })),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
})

export type CreateNoteRequest = z.infer<typeof CreateNoteRequestSchema>
export type NoteResponse = z.infer<typeof NoteResponseSchema>
```

### Handler (CQRSパターン)

外部からのエントリーポイント。コマンドとクエリを分離し、読み書きの責務を明確化。

```ts
// external/handler/note.query.server.ts
import 'server-only'
import { noteService } from '../service/note.service'
import { NoteResponseSchema } from '../dto/note.dto'

export async function listNotesServer(filters: NoteFilters) {
  // サービス層を通じてデータ取得
  const notes = await noteService.listNotes(filters)
  
  // DTOで型を保証して返す
  return notes.map(note => NoteResponseSchema.parse(note))
}
```

```ts
// external/handler/note.command.action.ts
'use server'
import { CreateNoteRequestSchema } from '../dto/note.dto'
import { noteService } from '../service/note.service'
import { getCurrentUser } from '@/shared/lib/auth'

export async function createNoteAction(request: unknown) {
  // 認証チェック
  const user = await getCurrentUser()
  if (!user) throw new Error('Unauthorized')

  // バリデーション
  const validated = CreateNoteRequestSchema.parse(request)
  
  // サービス層を通じて作成
  const note = await noteService.createNote(user.id, validated)
  
  return { id: note.id }
}
```

### Service (ビジネスロジック)

ビジネスロジックを実装し、**現在はDBに直接アクセスしているが、将来的には外部APIを呼び出すように変更可能**な設計。

```ts
// external/service/note.service.ts
import { apiClient } from '../client/api-client'
import type { CreateNoteRequest, NoteResponse } from '../dto/note.dto'

class NoteService {
  async createNote(
    userId: string, 
    input: CreateNoteRequest
  ): Promise<NoteResponse> {
    // 現在: Next.jsから直接DB操作（開発効率重視）
    const result = await db
      .insert(notes)
      .values({
        id: generateId(),
        title: input.title,
        templateId: input.templateId,
        ownerId: userId,
        status: 'Draft',
      })
      .returning()

    return this.formatNoteResponse(result[0])
  }

  // 将来: 外部APIを呼び出す実装に変更可能
  async createNoteViaAPI(
    userId: string, 
    input: CreateNoteRequest
  ): Promise<NoteResponse> {
    const response = await apiClient.post('/api/notes', {
      ...input,
      userId,
    })
    
    return NoteResponseSchema.parse(response.data)
  }

  async listNotes(filters: NoteFilters): Promise<NoteResponse[]> {
    // 現在はDB、将来はAPIエンドポイントを呼び出す
    // サービス層のインターフェースは変わらない
  }
}

export const noteService = new NoteService()
```

### Client (外部通信)

HTTPクライアントやDB接続を管理。将来的にはAPIクライアントに一本化。

```ts
// external/client/api-client.ts
import axios from 'axios'

// 将来的に外部APIサーバーと通信する際のクライアント
export const apiClient = axios.create({
  baseURL: process.env.API_BASE_URL || 'http://localhost:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 認証エラー処理
    }
    return Promise.reject(error)
  }
)
```

## データフロー

### 現在の実装（Next.js + DB直接接続）

```
Client Component
    ↓
Server Action (*.action.ts)
    ↓
Service (DB直接操作)
    ↓
Database (Drizzle ORM)
```

### 将来の実装（Next.js + 外部API）

```
Client Component
    ↓
Server Action (*.action.ts)
    ↓
Service (API呼び出し)
    ↓
External API (Go/Rust/etc)
    ↓
Database
```

## 移行戦略

1. **インターフェースの維持**: DTOの型定義は変更しない
2. **段階的移行**: サービス層のメソッドを1つずつAPIに置き換え
3. **フィーチャーフラグ**: 環境変数でDB直接/API切り替え

```ts
// external/service/base.service.ts
export abstract class BaseService {
  protected get useExternalAPI(): boolean {
    return process.env.USE_EXTERNAL_API === 'true'
  }

  protected async fetchData<T>(
    dbFetcher: () => Promise<T>,
    apiFetcher: () => Promise<T>
  ): Promise<T> {
    return this.useExternalAPI ? apiFetcher() : dbFetcher()
  }
}
```

## 命名規則

- **Query**: `*.query.server.ts` / `*.query.action.ts`
- **Command**: `*.command.server.ts` / `*.command.action.ts`
- **Server専用**: `import 'server-only'`を必ず記載
- **型定義**: DTOは入出力の契約、内部実装に依存しない

## ベストプラクティス

1. **DTOを変更しない**: バックエンドの実装が変わってもインターフェースは維持
2. **サービス層で抽象化**: DB操作とAPI呼び出しを同じインターフェースで扱う
3. **エラーハンドリング統一**: DB/APIどちらでも同じエラー型を返す
4. **テスト容易性**: サービス層をモックすることで、バックエンドに依存しないテストが可能