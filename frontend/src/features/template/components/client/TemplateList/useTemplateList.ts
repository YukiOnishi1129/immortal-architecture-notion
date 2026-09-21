"use client";

import { useSearchParams } from "next/navigation";
import { useTemplateListQuery } from "@/features/template/hooks/useTemplateQuery";
import type { TemplateFilters } from "@/features/template/types";

export function useTemplateList(initialFilters: TemplateFilters = {}) {
  const searchParams = useSearchParams();

  const filters: TemplateFilters = {
    q: searchParams.get("q") || initialFilters.q,
    page: Number(searchParams.get("page")) || initialFilters.page || 1,
    onlyMyTemplates:
      searchParams.get("onlyMyTemplates") === "true" ||
      initialFilters.onlyMyTemplates,
  };

  const { data: templates, isLoading } = useTemplateListQuery(filters);

  return {
    templates: templates || [],
    isLoading,
    filters,
  };
}
