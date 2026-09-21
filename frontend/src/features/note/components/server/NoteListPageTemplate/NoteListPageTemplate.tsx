import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { listNoteQuery } from "@/external/handler/note/note.query.server";
import { NoteListContainer } from "@/features/note/components/client/NoteList";
import { noteKeys } from "@/features/note/queries/keys";
import type { NoteStatus } from "@/features/note/types";
import { getQueryClient } from "@/shared/lib/query-client";

interface NoteListPageTemplateProps {
  status?: NoteStatus;
  q?: string;
  page?: number;
  templateId?: string;
}

export async function NoteListPageTemplate({
  q,
  page,
  templateId,
}: NoteListPageTemplateProps) {
  const queryClient = getQueryClient();

  // 公開ノートのみを強制的に表示
  const filters = {
    status: "Publish" as const, // 常に公開済みのみ
    q,
    page: page || 1,
    templateId,
  };

  // データをプリフェッチ（公開済みノートのみ）
  await queryClient.prefetchQuery({
    queryKey: noteKeys.list(filters),
    queryFn: () => listNoteQuery(filters),
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <NoteListContainer initialFilters={filters} />
    </HydrationBoundary>
  );
}
