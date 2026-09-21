import { MyNoteListPageTemplate } from "@/features/note/components/server/MyNoteListPageTemplate";
import type { NoteStatus } from "@/features/note/types";

export default async function MyNotesPage({
  searchParams,
}: PageProps<"/my-notes">) {
  const params = await searchParams;

  const status = params.status as NoteStatus | undefined;
  const q = typeof params.q === "string" ? params.q : undefined;
  const page =
    typeof params.page === "string" ? Number(params.page) : undefined;

  return <MyNoteListPageTemplate status={status} q={q} page={page} />;
}
