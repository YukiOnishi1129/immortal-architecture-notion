import { useQuery } from "@tanstack/react-query";
import {
  getTemplateByIdQueryAction,
  listTemplatesQueryAction,
} from "@/external/handler/template/template.query.action";
import { templateKeys } from "@/features/template/queries/keys";
import type { TemplateFilters } from "@/features/template/types";

export function useTemplateListQuery(filters: TemplateFilters) {
  return useQuery({
    queryKey: templateKeys.list(filters),
    queryFn: () => listTemplatesQueryAction(filters),
  });
}

export function useTemplateQuery(templateId: string) {
  return useQuery({
    queryKey: templateKeys.detail(templateId),
    queryFn: () => getTemplateByIdQueryAction({ id: templateId }),
  });
}
