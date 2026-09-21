import { MyNoteDetailPageTemplate } from "@/features/note/components/server/MyNoteDetailPageTemplate";

export default async function MyNoteDetailPage({
  params,
}: PageProps<"/my-notes/[id]">) {
  const { id } = await params;
  return <MyNoteDetailPageTemplate noteId={id} />;
}
