"use client";

import { useSearchParams } from "next/navigation";
import { NOTE_STATUS } from "@/features/note/constants";
import { useNoteListQuery } from "@/features/note/hooks/useNoteListQuery";
import type { NoteFilters } from "@/features/note/types";

export function useNoteList(initialFilters: NoteFilters = {}) {
  const searchParams = useSearchParams();

  const pageParam = searchParams.get("page");

  // 公開ノートのみを強制的に表示
  const filters: NoteFilters = {
    q: searchParams.get("q") || initialFilters.q,
    status: NOTE_STATUS.PUBLISH, // 常に公開のみ表示
    page: pageParam ? Number.parseInt(pageParam, 10) : initialFilters.page || 1,
  };

  const { data: notes, isLoading } = useNoteListQuery(filters);

  return {
    notes: notes || [],
    isLoading,
    filters,
  };
}
