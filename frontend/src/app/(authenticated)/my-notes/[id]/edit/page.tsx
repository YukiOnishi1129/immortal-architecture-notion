import { MyNoteEditPageTemplate } from "@/features/note/components/server/MyNoteEditPageTemplate";

export default async function MyNoteEditPage({
  params,
}: PageProps<"/my-notes/[id]/edit">) {
  const { id } = await params;

  return <MyNoteEditPageTemplate noteId={id} />;
}
