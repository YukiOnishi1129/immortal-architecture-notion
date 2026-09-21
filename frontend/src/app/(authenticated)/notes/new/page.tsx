import { NoteNewPageTemplate } from "@/features/note/components/server/NoteNewPageTemplate";

export default async function CreateNotePage({
  searchParams,
}: PageProps<"/notes/new">) {
  const params = await searchParams;
  const templateId =
    typeof params.templateId === "string" ? params.templateId : undefined;

  return <NoteNewPageTemplate initialTemplateId={templateId} />;
}
