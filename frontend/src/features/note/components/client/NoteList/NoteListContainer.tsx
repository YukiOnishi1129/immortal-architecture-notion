"use client";

import type { NoteFilters } from "@/features/note/types";
import { NoteListPresenter } from "./NoteListPresenter";
import { useNoteList } from "./useNoteList";

interface NoteListContainerProps {
  initialFilters?: NoteFilters;
}

export function NoteListContainer({
  initialFilters = {},
}: NoteListContainerProps) {
  const { notes, isLoading, filters } = useNoteList(initialFilters);

  return (
    <NoteListPresenter notes={notes} isLoading={isLoading} filters={filters} />
  );
}
