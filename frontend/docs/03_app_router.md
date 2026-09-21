# App Router 設計ガイド

## 基本方針

- `page.tsx`と`layout.tsx`は全てRSC (React Server Component)
- `error.tsx`のみClient Component
- ビジネスロジックは`features/`に委譲
- ルート構造で認証状態を表現
- Next.js 15+のグローバル型定義（`LayoutProps`/`PageProps`）を活用

## ルートグループ戦略

### 認証別グループ

```
app/
├─ (guest)/          # 未ログインユーザー向け
│  ├─ login/
│  └─ signup/
├─ (authenticated)/  # ログイン必須
│  ├─ notes/
│  ├─ templates/
│  └─ me/
└─ (neutral)/        # 認証不問
   ├─ terms/
   └─ privacy/
```

### グループ別設定

| グループ | Layout | 認証チェック | 共通UI |
|---------|--------|------------|---------|
| `(guest)` | シンプル | リダイレクト | なし |
| `(authenticated)` | フル機能 | 必須 | Header, Sidebar |
| `(neutral)` | 最小限 | なし | Footer のみ |

## ページコンポーネントパターン

### 基本構造

```tsx
// app/notes/[noteId]/page.tsx
import { NoteDetailPageTemplate } from '@/features/note/components/server/NoteDetailPageTemplate'

// Next.js 15+のグローバル型定義を使用（importなし）
export default async function NoteDetailPage(props: PageProps<'/notes/[noteId]'>) {
  const params = await props.params
  const searchParams = await props.searchParams
  
  return <NoteDetailPageTemplate 
    noteId={params.noteId}
    searchParams={searchParams} 
  />
}
```

### メタデータ設定

```tsx
// app/notes/layout.tsx
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'ノート一覧 | Mini Notion',
  description: '設計メモを構造化して管理',
}

// Next.js 15+のグローバル型定義を使用（importなし）
export default function NotesLayout(props: LayoutProps<'/notes'>) {
  return <>{props.children}</>
}
```

## 認証レイアウト実装

```tsx
// app/(authenticated)/layout.tsx
import { redirect } from 'next/navigation'
import { getServerSession } from 'next-auth'
import { authOptions } from '@/shared/lib/auth'
import { AuthenticatedLayoutWrapper } from '@/shared/components/layout/server/AuthenticatedLayoutWrapper'

export default async function AuthenticatedLayout(props: LayoutProps<'/'>) {
  const session = await getServerSession(authOptions)
  
  if (!session) {
    redirect('/login')
  }

  return (
    <AuthenticatedLayoutWrapper user={session.user}>
      {props.children}
    </AuthenticatedLayoutWrapper>
  )
}
```

## エラーハンドリング

```tsx
// app/(authenticated)/error.tsx
'use client'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px]">
      <h2 className="text-2xl font-bold mb-4">エラーが発生しました</h2>
      <p className="text-gray-600 mb-6">{error.message}</p>
      <button
        onClick={reset}
        className="px-4 py-2 bg-primary text-white rounded-md"
      >
        再試行
      </button>
    </div>
  )
}
```

## ローディング状態

```tsx
// app/(authenticated)/notes/loading.tsx
export default function Loading() {
  return (
    <div className="animate-pulse">
      <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
      <div className="space-y-3">
        <div className="h-4 bg-gray-200 rounded"></div>
        <div className="h-4 bg-gray-200 rounded"></div>
      </div>
    </div>
  )
}
```