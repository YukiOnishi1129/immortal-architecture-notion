import { NoteListPageTemplate } from "@/features/note/components/server/NoteListPageTemplate";
import type { NoteStatus } from "@/features/note/types";

export default async function NotesPage({ searchParams }: PageProps<"/notes">) {
  const params = await searchParams;

  const status = params.status as NoteStatus | undefined;
  const q = typeof params.q === "string" ? params.q : undefined;
  const page =
    typeof params.page === "string" ? Number(params.page) : undefined;
  const templateId =
    typeof params.templateId === "string" ? params.templateId : undefined;

  return (
    <NoteListPageTemplate
      status={status}
      q={q}
      page={page}
      templateId={templateId}
    />
  );
}
