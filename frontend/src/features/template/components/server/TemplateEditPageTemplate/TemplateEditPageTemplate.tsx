import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { notFound } from "next/navigation";
import { getTemplateByIdQuery } from "@/external/handler/template/template.query.server";
import { TemplateEditForm } from "@/features/template/components/client/TemplateEditForm";
import { templateKeys } from "@/features/template/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

interface TemplateEditPageTemplateProps {
  templateId: string;
}

export async function TemplateEditPageTemplate({
  templateId,
}: TemplateEditPageTemplateProps) {
  const template = await getTemplateByIdQuery({ id: templateId });

  if (!template) {
    notFound();
  }

  const queryClient = getQueryClient();
  await queryClient.prefetchQuery({
    queryKey: templateKeys.detail(templateId),
    queryFn: () => template,
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <TemplateEditForm templateId={templateId} />
      </div>
    </HydrationBoundary>
  );
}
