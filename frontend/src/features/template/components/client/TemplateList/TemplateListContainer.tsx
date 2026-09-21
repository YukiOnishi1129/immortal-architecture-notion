"use client";

import type { TemplateFilters } from "@/features/template/types";
import { TemplateListPresenter } from "./TemplateListPresenter";
import { useTemplateList } from "./useTemplateList";

interface TemplateListContainerProps {
  initialFilters?: TemplateFilters;
}

export function TemplateListContainer({
  initialFilters = {},
}: TemplateListContainerProps) {
  const { templates, isLoading } = useTemplateList(initialFilters);

  return <TemplateListPresenter templates={templates} isLoading={isLoading} />;
}
