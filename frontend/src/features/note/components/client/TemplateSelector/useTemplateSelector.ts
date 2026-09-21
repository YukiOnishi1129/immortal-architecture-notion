"use client";

import { useQuery } from "@tanstack/react-query";
import { listTemplatesQueryAction } from "@/external/handler/template/template.query.action";
import { templateKeys } from "@/features/template/queries/keys";

export function useTemplateSelector() {
  const { data, isLoading } = useQuery({
    queryKey: templateKeys.list({}),
    queryFn: () => listTemplatesQueryAction(),
  });

  return {
    templates: data ?? [],
    isLoading,
  };
}
