"use client";

import { TemplateListFilterPresenter } from "./TemplateListFilterPresenter";
import { useTemplateListFilter } from "./useTemplateListFilter";

export function TemplateListFilterContainer() {
  const {
    searchQuery,
    isPending,
    onlyMyTemplates,
    setSearchQuery,
    handleSearch,
    handleOnlyMyTemplatesChange,
  } = useTemplateListFilter();

  return (
    <TemplateListFilterPresenter
      searchQuery={searchQuery}
      isPending={isPending}
      onlyMyTemplates={onlyMyTemplates}
      onSearchQueryChange={setSearchQuery}
      onSearchSubmit={handleSearch}
      onOnlyMyTemplatesChange={handleOnlyMyTemplatesChange}
    />
  );
}
