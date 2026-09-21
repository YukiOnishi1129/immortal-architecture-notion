import { TemplateListPageTemplate } from "@/features/template/components/server/TemplateListPageTemplate";

export default function TemplatesPage({
  searchParams,
}: PageProps<"/templates">) {
  return <TemplateListPageTemplate searchParams={searchParams} />;
}
