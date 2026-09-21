import { notFound, redirect } from "next/navigation";
import { getTemplateByIdQuery } from "@/external/handler/template/template.query.server";
import { getSessionServer } from "@/features/auth/servers/auth.server";
import { TemplateEditPageTemplate } from "@/features/template/components/server/TemplateEditPageTemplate";

export default async function TemplateEditPage({
  params,
}: PageProps<"/templates/[id]/edit">) {
  const { id } = await params;
  const [template, session] = await Promise.all([
    getTemplateByIdQuery({ id }),
    getSessionServer(),
  ]);

  if (!template) {
    notFound();
  }

  // Check if user is owner
  if (session?.account?.id !== template.ownerId) {
    notFound(); // or redirect to 403
  }

  // Check if template is used by notes
  if (template.isUsed) {
    redirect(`/templates/${id}`);
  }

  return <TemplateEditPageTemplate templateId={id} />;
}
