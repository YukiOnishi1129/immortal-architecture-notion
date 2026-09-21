"use client";

import type { NoteFilters } from "@/features/note/types";
import { MyNoteListPresenter } from "./MyNoteListPresenter";
import { useMyNoteList } from "./useMyNoteList";

type MyNotesContainerProps = {
  initialFilters: NoteFilters;
};

export function MyNoteListContainer({ initialFilters }: MyNotesContainerProps) {
  const { notes, isLoading, filters } = useMyNoteList({ initialFilters });

  return (
    <MyNoteListPresenter
      notes={notes}
      isLoading={isLoading}
      filters={filters}
    />
  );
}
