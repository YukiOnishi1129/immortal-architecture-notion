# Features ディレクトリ設計

## 概要

Featuresディレクトリは、アプリケーションの機能を**ドメイン単位**で整理します。各機能は独立したモジュールとして設計され、高い凝集性と低い結合性を保ちます。

## ディレクトリ構造

```
features/
├─ notes/        # ノート機能
├─ templates/    # テンプレート機能
├─ auth/         # 認証機能
└─ account/      # アカウント機能
```

## 機能モジュールの内部構造

```
features/note/
├─ components/
│  ├─ server/    # Server Components
│  │  ├─ NoteListPageTemplate/
│  │  │  ├─ index.ts
│  │  │  └─ NoteListPageTemplate.tsx
│  │  └─ NoteDetailPageTemplate/
│  │     ├─ index.ts
│  │     └─ NoteDetailPageTemplate.tsx
│  └─ client/    # Client Components
│     ├─ NoteList/
│     │  ├─ index.ts
│     │  ├─ NoteListContainer.tsx
│     │  ├─ NoteListPresenter.tsx
│     │  └─ useNoteList.ts
│     └─ NoteForm/
├─ hooks/        # カスタムフック
│  ├─ useNoteQuery.ts
│  └─ useNoteMutation.ts
├─ queries/      # TanStack Query関連
│  ├─ keys.ts
│  └─ helpers.ts
├─ actions/      # Server Actions
│  ├─ createNote.ts
│  └─ updateNote.ts
├─ types/        # 型定義
│  └─ index.ts
└─ utils/        # ユーティリティ
   └─ validation.ts
```

## Container/Presenterパターン

### Container (ロジック層)

**重要な制約: ContainerはDOMを直接レンダリングせず、対応するPresenterにpropsを渡すだけにすること。**

Containerの責務:
- カスタムフックを使ってデータを取得する
- イベントハンドラーを定義する
- **Presenterコンポーネントをレンダリングしてpropsを渡す**
- **DOM要素（div、button、linkなど）を直接レンダリングしない**

```tsx
// features/note/components/client/NoteList/NoteListContainer.tsx
'use client'

import { NoteListPresenter } from './NoteListPresenter'
import { useNoteList } from './useNoteList'

interface NoteListContainerProps {
  initialFilters?: NoteFilters
}

export function NoteListContainer({ initialFilters }: NoteListContainerProps) {
  const {
    notes,
    isLoading,
    filters,
    updateFilters,
    handleDelete,
  } = useNoteList(initialFilters)

  // ✅ PresenterにpropsだけをRenderingする
  return (
    <NoteListPresenter
      notes={notes}
      isLoading={isLoading}
      filters={filters}
      onFilterChange={updateFilters}
      onDelete={handleDelete}
    />
  )
}
```

悪い例 ❌:

```tsx
// ❌ ContainerでDOMを直接レンダリングしている
export function NoteListContainer({ initialFilters }: NoteListContainerProps) {
  const { notes, isLoading, filters } = useNoteList(initialFilters)

  return (
    <div className="space-y-6">  {/* ❌ ContainerでDOM要素を書いている */}
      <div className="bg-white p-6">
        <h1>タイトル</h1>
        <FilterBar filters={filters} />
      </div>
      <NoteListPresenter notes={notes} isLoading={isLoading} />
    </div>
  )
}
```

このような場合は、全てのDOMをPresenterに移動すること。

### Presenter (表示層)

```tsx
// features/note/components/client/NoteList/NoteListPresenter.tsx
import { NoteCard } from '../NoteCard'
import { FilterBar } from '../FilterBar'
import { LoadingSpinner } from '@/shared/components/ui/LoadingSpinner'

interface NoteListPresenterProps {
  notes: Note[]
  isLoading: boolean
  filters: NoteFilters
  onFilterChange: (filters: NoteFilters) => void
  onDelete: (noteId: string) => void
}

export function NoteListPresenter({
  notes,
  isLoading,
  filters,
  onFilterChange,
  onDelete,
}: NoteListPresenterProps) {
  if (isLoading) return <LoadingSpinner />

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} onChange={onFilterChange} />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {notes.map((note) => (
          <NoteCard
            key={note.id}
            note={note}
            onDelete={() => onDelete(note.id)}
          />
        ))}
      </div>
    </div>
  )
}
```

### カスタムフック

```tsx
// features/note/components/client/NoteList/useNoteList.ts
import { useState, useCallback } from 'react'
import { useNoteListQuery } from '@/features/note/hooks/useNoteQuery'
import { useDeleteNoteMutation } from '@/features/note/hooks/useNoteMutation'

export function useNoteList(initialFilters?: NoteFilters) {
  const [filters, setFilters] = useState(initialFilters || {})
  const { data, isLoading } = useNoteListQuery(filters)
  const deleteMutation = useDeleteNoteMutation()

  // イベントハンドラーはuseCallbackで最適化
  const handleDelete = useCallback(async (noteId: string) => {
    await deleteMutation.mutateAsync(noteId)
  }, [deleteMutation])

  const updateFilters = useCallback((newFilters: NoteFilters) => {
    setFilters(newFilters)
  }, [])

  return {
    notes: data?.notes || [],
    isLoading,
    filters,
    updateFilters,
    handleDelete,
  }
}
```

## Server Componentsテンプレート

Server Components は専用のディレクトリを作成し、index.tsでエクスポートを管理します。

### ディレクトリ構成

```
server/
├─ LoginPageTemplate/
│  ├─ index.ts                 # export { LoginPageTemplate } from './LoginPageTemplate'
│  └─ LoginPageTemplate.tsx    # 実際のコンポーネント実装
└─ NoteListPageTemplate/
   ├─ index.ts
   └─ NoteListPageTemplate.tsx
```

### 実装例

```tsx
// features/note/components/server/NoteListPageTemplate/NoteListPageTemplate.tsx
import { HydrationBoundary, dehydrate } from '@tanstack/react-query'
import { getQueryClient } from '@/shared/lib/query-client'
import { noteKeys } from '@/features/note/queries/keys'
import { listNotesServer } from '@/external/handler/note.query.server'
import { NoteListContainer } from '../../client/NoteList'

interface NoteListPageTemplateProps {
  searchParams: { [key: string]: string | string[] | undefined }
}

export async function NoteListPageTemplate({ searchParams }: NoteListPageTemplateProps) {
  const queryClient = getQueryClient()
  const filters = { 
    status: searchParams.status as NoteStatus,
    q: searchParams.q as string 
  }

  await queryClient.prefetchQuery({
    queryKey: noteKeys.list(filters),
    queryFn: () => listNotesServer(filters),
  })

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <NoteListContainer initialFilters={filters} />
    </HydrationBoundary>
  )
}
```

```tsx
// features/note/components/server/NoteListPageTemplate/index.ts
export { NoteListPageTemplate } from './NoteListPageTemplate'
```

## Client Componentsの命名規則

index.tsでエクスポートする際は、より具体的で意味のある名前に変更します：

```tsx
// features/auth/components/client/Login/index.ts
export { LoginContainer as LoginForm } from './LoginContainer'

// features/note/components/client/NoteList/index.ts  
export { NoteListContainer as NoteList } from './NoteListContainer'
```

## Presenterコンポーネントの使用ルール

### 重要な制約

**Presenterコンポーネントは、同じ機能ディレクトリ内のContainerからのみ呼び出すこと。**

他の機能ディレクトリや異なるコンポーネントから直接Presenterを呼び出すことは禁止です。

### 良い例 ✅

```tsx
// features/note/components/client/MyNoteList/MyNoteListContainer.tsx
import { MyNoteListPresenter } from './MyNoteListPresenter'

export function MyNoteListContainer({ initialFilters }: Props) {
  const { notes, isLoading, filters } = useMyNoteList({ initialFilters })

  // 同じディレクトリ内のPresenterを呼び出す
  return (
    <MyNoteListPresenter
      notes={notes}
      isLoading={isLoading}
      filters={filters}
    />
  )
}
```

### 悪い例 ❌

```tsx
// features/note/components/client/MyNoteList/MyNoteListPresenter.tsx
import { NoteListPresenter } from '../NoteList/NoteListPresenter' // ❌ 別のPresenterを呼び出している

export function MyNoteListPresenter({ notes, isLoading }: Props) {
  return (
    <div>
      <h1>マイノート</h1>
      {/* ❌ 他のコンポーネントのPresenterを直接呼び出すのは禁止 */}
      <NoteListPresenter notes={notes} isLoading={isLoading} />
    </div>
  )
}
```

### 正しい対処法

他のコンポーネントの表示ロジックを再利用したい場合は、以下のいずれかの方法を取ります：

1. **Presenterの実装を直接コピーして独自に実装する**
2. **共通部分をshared/componentsに切り出す**
3. **Containerを呼び出す（Presenterではなく）**

```tsx
// features/note/components/client/MyNoteList/MyNoteListPresenter.tsx
import { NoteList } from '../NoteList' // ✅ Containerを呼び出す（index.tsでエクスポートされたもの）

export function MyNoteListPresenter({ filters }: Props) {
  return (
    <div>
      <h1>マイノート</h1>
      {/* ✅ Containerを呼び出すのはOK */}
      <NoteList initialFilters={filters} />
    </div>
  )
}
```

ただし、この場合はMyNoteListPresenterが独自のロジックと表示を持つべきなので、通常は方法1（独自に実装）を選択します。

## コンポーネント分割のルール

### 1ファイル1コンポーネントの原則

**すべてのClient Componentは1ファイルにつき1コンポーネントのみ定義すること。**

複数のコンポーネントが1ファイルに存在する場合は、以下のルールに従って分割します：

#### View専用コンポーネント（ロジックなし）の場合

**同じディレクトリ内に配置**します。

```
NoteList/
├─ index.ts
├─ NoteListContainer.tsx      # メインのContainer
├─ NoteListPresenter.tsx       # メインのPresenter
├─ NoteListItem.tsx           # ✅ View専用の子コンポーネント（同じディレクトリ）
└─ NoteListSkeleton.tsx       # ✅ View専用の子コンポーネント（同じディレクトリ）
```

**例：View専用コンポーネント**
```tsx
// NoteListItem.tsx
import type { Note } from '@/features/note/types'

interface NoteListItemProps {
  note: Note
  onSelect: (id: string) => void
}

export function NoteListItem({ note, onSelect }: NoteListItemProps) {
  return (
    <div onClick={() => onSelect(note.id)}>
      <h3>{note.title}</h3>
      <p>{note.status}</p>
    </div>
  )
}
```

#### ロジックを含むコンポーネントの場合

**client配下に新しいディレクトリを作成**します。

```
client/
├─ NoteList/
│  ├─ index.ts
│  ├─ NoteListContainer.tsx
│  ├─ NoteListPresenter.tsx
│  └─ useNoteList.ts
└─ NoteListFilter/              # ✅ ロジックを含むため別ディレクトリ
   ├─ index.ts
   ├─ NoteListFilterContainer.tsx
   ├─ NoteListFilterPresenter.tsx
   └─ useNoteListFilter.ts
```

### Presenterのルール

**Presenterコンポーネントはロジックを持たず、propsで渡されたデータを表示するのみ。**

#### ❌ 悪い例：Presenterにロジックがある

```tsx
export function NoteListPresenter({ notes }: Props) {
  // ❌ Presenterでフィルタリングロジックを持っている
  const [filter, setFilter] = useState('all')
  const filteredNotes = notes.filter(note =>
    filter === 'all' ? true : note.status === filter
  )

  return (
    <div>
      <select onChange={(e) => setFilter(e.target.value)}>
        <option value="all">すべて</option>
        <option value="Draft">下書き</option>
      </select>
      {filteredNotes.map(note => <NoteCard key={note.id} note={note} />)}
    </div>
  )
}
```

#### ✅ 良い例：ロジックはContainerとHookに分離

```tsx
// NoteListPresenter.tsx
export function NoteListPresenter({
  notes,
  filter,
  onFilterChange
}: Props) {
  // ✅ ロジックなし、propsで渡されたものを表示するのみ
  return (
    <div>
      <select value={filter} onChange={(e) => onFilterChange(e.target.value)}>
        <option value="all">すべて</option>
        <option value="Draft">下書き</option>
      </select>
      {notes.map(note => <NoteCard key={note.id} note={note} />)}
    </div>
  )
}
```

```tsx
// NoteListContainer.tsx
export function NoteListContainer({ initialNotes }: Props) {
  // ✅ ロジックはContainerとHookに集約
  const { notes, filter, handleFilterChange } = useNoteList(initialNotes)

  return (
    <NoteListPresenter
      notes={notes}
      filter={filter}
      onFilterChange={handleFilterChange}
    />
  )
}
```

## ベストプラクティス

1. **1ファイル1コンポーネント**: 複数のコンポーネントがある場合は必ず分割する
2. **Presenterは純粋な表示のみ**: ロジック（useState、useEffect等）を持たない
3. **ロジックの配置**: Container + Custom Hookにロジックを集約
4. **View専用コンポーネントの配置**: ロジックがないなら同じディレクトリ、ロジックがあるなら別ディレクトリ
5. **単一責任の原則**: 各コンポーネントは1つの責任のみを持つ
6. **再利用性**: 汎用的なコンポーネントは`shared/`へ移動
7. **テスタビリティ**: PresenterはPropsのみに依存
8. **型安全性**: 全てのインターフェースを明示的に定義
9. **Presenterの独立性**: Presenterは他のPresenterを呼び出さない（同じディレクトリ内のContainerからのみ呼び出される）
10. **命名規則**:
    - ファイル名とコンポーネント名を一致させる（アッパーキャメルケース）
    - Server ComponentsはxxxPageTemplateの命名規則
    - Client Componentsはindex.tsで適切な名前でエクスポート