import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { notFound } from "next/navigation";
import { getTemplateByIdQuery } from "@/external/handler/template/template.query.server";
import { getSessionServer } from "@/features/auth/servers/auth.server";
import { TemplateDetail } from "@/features/template/components/client/TemplateDetail";
import { templateKeys } from "@/features/template/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

interface TemplateDetailPageTemplateProps {
  templateId: string;
}

export async function TemplateDetailPageTemplate({
  templateId,
}: TemplateDetailPageTemplateProps) {
  const [template, session] = await Promise.all([
    getTemplateByIdQuery({ id: templateId }),
    getSessionServer(),
  ]);

  if (!template) {
    notFound();
  }

  const queryClient = getQueryClient();
  await queryClient.prefetchQuery({
    queryKey: templateKeys.detail(templateId),
    queryFn: () => template,
  });

  const isOwner = session?.account?.id === template.ownerId;

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <TemplateDetail templateId={templateId} isOwner={isOwner} />
      </div>
    </HydrationBoundary>
  );
}
