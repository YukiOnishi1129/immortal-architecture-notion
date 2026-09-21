# TanStack Query 実装ガイド

## 概要

TanStack Query を使用してサーバー状態を管理し、Next.js App Router の Server Components と連携させます。

## セットアップ

### Provider 設定

```tsx
// shared/providers/query-provider.tsx
"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useState } from "react";

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 0, // RSCのhydrateデータを常に優先
            gcTime: 5 * 60 * 1000, // 5分（デフォルト）
            refetchOnWindowFocus: false,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

### staleTime と gcTime の設定意図

| 設定      | 値   | 理由                                                                     |
| --------- | ---- | ------------------------------------------------------------------------ |
| staleTime | 0    | RSC で hydrate されたデータを常に優先。ページ遷移時に最新データを反映    |
| gcTime    | 5 分 | キャッシュをメモリに保持。楽観的更新やコンポーネント間のデータ共有で使用 |

**staleTime: 0 でも TanStack Query を使う意味:**

- 同一ページ内での状態共有（サイドバーとメインコンテンツなど）
- Mutation 後の楽観的更新
- ローディング状態の管理

### サーバー用 QueryClient

```tsx
// shared/lib/query-client.ts
import { QueryClient } from "@tanstack/react-query";
import { cache } from "react";

export const getQueryClient = cache(() => new QueryClient());
```

## クエリキーの管理

```ts
// features/note/queries/keys.ts
export const noteKeys = {
  all: ["notes"] as const,
  lists: () => [...noteKeys.all, "list"] as const,
  list: (filters: NoteFilters) => [...noteKeys.lists(), filters] as const,
  details: () => [...noteKeys.all, "detail"] as const,
  detail: (id: string) => [...noteKeys.details(), id] as const,
};
```

## サーバーサイドプリフェッチ

```tsx
// features/note/components/server/NotesPageTemplate.tsx
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { getQueryClient } from "@/shared/lib/query-client";
import { noteKeys } from "@/features/note/queries/keys";
import { listNotesServer } from "@/external/handler/note.query.server";
import { NoteListContainer } from "../client/NoteList";

interface NotesPageTemplateProps {
  filters?: NoteFilters;
}

export async function NotesPageTemplate({
  filters = {},
}: NotesPageTemplateProps) {
  const queryClient = getQueryClient();

  // データをプリフェッチ
  await queryClient.prefetchQuery({
    queryKey: noteKeys.list(filters),
    queryFn: () => listNotesServer(filters),
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <NoteListContainer initialFilters={filters} />
    </HydrationBoundary>
  );
}
```

## クライアントサイド Query

### Server Action の使用

API Routes ではなく Server Actions を使用してデータフェッチを行います。これにより型安全性が向上し、エンドツーエンドの型推論が可能になります。

```ts
// features/note/hooks/useNoteQuery.ts
"use client";

import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { noteKeys } from "../queries/keys";
import {
  listNotesAction,
  getNoteDetailAction,
} from "@/external/handler/note.query.action";

// Server Actionを直接使用
export function useNoteListQuery(filters: NoteFilters) {
  return useQuery({
    queryKey: noteKeys.list(filters),
    queryFn: () => listNotesAction(filters), // Server Action呼び出し
  });
}

export function useNoteDetailQuery(noteId: string) {
  return useSuspenseQuery({
    queryKey: noteKeys.detail(noteId),
    queryFn: () => getNoteDetailAction(noteId), // Server Action呼び出し
  });
}
```

**重要**: API Routes (`/api/notes`) ではなく、`external/handler` ディレクトリの Server Actions を使用してください。

## Mutation 実装

```ts
// features/note/hooks/useNoteMutation.ts
"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "@/shared/hooks/use-toast";
import { noteKeys } from "../queries/keys";
import {
  createNoteAction,
  updateNoteAction,
  deleteNoteAction,
} from "@/external/handler/note.command.action";

export function useCreateNoteMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createNoteAction,
    onSuccess: (data) => {
      // 関連するクエリを無効化
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
      toast({ title: "ノートを作成しました" });
    },
    onError: (error) => {
      toast({
        title: "エラーが発生しました",
        description: error.message,
        variant: "destructive",
      });
    },
  });
}

export function useUpdateNoteMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateNoteAction,
    onSuccess: (_, variables) => {
      // 特定のノートと一覧を無効化
      queryClient.invalidateQueries({
        queryKey: noteKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({
        queryKey: noteKeys.lists(),
      });
    },
  });
}
```

## 楽観的更新

```ts
export function useUpdateNoteOptimistic() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateNoteAction,
    onMutate: async (newNote) => {
      // 既存のクエリをキャンセル
      await queryClient.cancelQueries({
        queryKey: noteKeys.detail(newNote.id),
      });

      // 現在のデータを保存
      const previousNote = queryClient.getQueryData(
        noteKeys.detail(newNote.id)
      );

      // 楽観的に更新
      queryClient.setQueryData(noteKeys.detail(newNote.id), newNote);

      return { previousNote };
    },
    onError: (err, newNote, context) => {
      // エラー時は元に戻す
      if (context?.previousNote) {
        queryClient.setQueryData(
          noteKeys.detail(newNote.id),
          context.previousNote
        );
      }
    },
    onSettled: (_, __, variables) => {
      // 最終的にサーバーデータで同期
      queryClient.invalidateQueries({
        queryKey: noteKeys.detail(variables.id),
      });
    },
  });
}
```

## 無限スクロール

```ts
// features/note/hooks/useInfiniteNotes.ts
import { useInfiniteQuery } from "@tanstack/react-query";

export function useInfiniteNotes(filters: NoteFilters) {
  return useInfiniteQuery({
    queryKey: noteKeys.list(filters),
    queryFn: ({ pageParam }) =>
      listNotesAction({ ...filters, page: pageParam }),
    getNextPageParam: (lastPage) => lastPage.nextPage,
    initialPageParam: 1,
  });
}
```

## パフォーマンス最適化

### Suspense との統合

```tsx
// features/note/components/client/NoteDetail.tsx
import { Suspense } from "react";
import { ErrorBoundary } from "react-error-boundary";

export function NoteDetail({ noteId }: { noteId: string }) {
  return (
    <ErrorBoundary fallback={<ErrorFallback />}>
      <Suspense fallback={<NoteSkeleton />}>
        <NoteDetailContent noteId={noteId} />
      </Suspense>
    </ErrorBoundary>
  );
}

function NoteDetailContent({ noteId }: { noteId: string }) {
  // useSuspenseQueryを使用
  const { data: note } = useNoteDetailQuery(noteId);

  return <NotePresenter note={note} />;
}
```

### 選択的な無効化

```ts
// 影響範囲を限定した無効化
onSuccess: async (_, { noteId, templateId }) => {
  await Promise.all([
    // 特定のノートのみ
    queryClient.invalidateQueries({
      queryKey: noteKeys.detail(noteId),
      exact: true,
    }),
    // 同じテンプレートのノート一覧
    queryClient.invalidateQueries({
      queryKey: noteKeys.list({ templateId }),
    }),
  ]);
};
```
