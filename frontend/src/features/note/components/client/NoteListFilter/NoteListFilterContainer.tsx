"use client";

import type { NoteFilters } from "@/features/note/types";
import { NoteListFilterPresenter } from "./NoteListFilterPresenter";
import { useNoteListFilter } from "./useNoteListFilter";

interface NoteListFilterContainerProps {
  filters: NoteFilters;
}

export function NoteListFilterContainer({
  filters,
}: NoteListFilterContainerProps) {
  const {
    searchQuery,
    statusFilter,
    isPending,
    setSearchQuery,
    handleSearch,
    handleStatusChange,
  } = useNoteListFilter(filters);

  return (
    <NoteListFilterPresenter
      searchQuery={searchQuery}
      statusFilter={statusFilter}
      isPending={isPending}
      onSearchQueryChange={setSearchQuery}
      onSearchSubmit={handleSearch}
      onStatusChange={handleStatusChange}
    />
  );
}
