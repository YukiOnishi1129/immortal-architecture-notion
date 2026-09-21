import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { getTemplateByIdQuery } from "@/external/handler/template/template.query.server";
import { NoteNewForm } from "@/features/note/components/client/NoteNewForm";
import { templateKeys } from "@/features/template/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

interface NoteNewPageTemplateProps {
  initialTemplateId?: string;
}

export async function NoteNewPageTemplate({
  initialTemplateId,
}: NoteNewPageTemplateProps) {
  const queryClient = getQueryClient();

  // initialTemplateIdがある場合、テンプレートをプリフェッチ
  if (initialTemplateId) {
    await queryClient.prefetchQuery({
      queryKey: templateKeys.detail(initialTemplateId),
      queryFn: () => getTemplateByIdQuery({ id: initialTemplateId }),
    });
  }

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container mx-auto py-6">
        <NoteNewForm initialTemplateId={initialTemplateId} />
      </div>
    </HydrationBoundary>
  );
}
