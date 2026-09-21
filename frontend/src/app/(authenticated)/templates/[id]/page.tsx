import { TemplateDetailPageTemplate } from "@/features/template/components/server/TemplateDetailPageTemplate";

export default async function TemplateDetailPage({
  params,
}: PageProps<"/templates/[id]">) {
  const { id } = await params;

  return <TemplateDetailPageTemplate templateId={id} />;
}
