"use client";

import { PublicNoteListFilterPresenter } from "./PublicNoteListFilterPresenter";
import { usePublicNoteListFilter } from "./usePublicNoteListFilter";

interface PublicNoteListFilterContainerProps {
  filters?: {
    q?: string;
  };
}

export function PublicNoteListFilterContainer({
  filters,
}: PublicNoteListFilterContainerProps) {
  const { searchQuery, isPending, setSearchQuery, handleSearch } =
    usePublicNoteListFilter(filters);

  return (
    <PublicNoteListFilterPresenter
      searchQuery={searchQuery}
      isPending={isPending}
      onSearchQueryChange={setSearchQuery}
      onSearchSubmit={handleSearch}
    />
  );
}
