import { NoteDetailPageTemplate } from "@/features/note/components/server/NoteDetailPageTemplate";

export default async function NoteDetailPage({
  params,
}: PageProps<"/notes/[id]">) {
  const { id } = await params;
  return <NoteDetailPageTemplate noteId={id} />;
}
